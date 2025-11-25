package outbound

//go:generate go run github.com/v2fly/v2ray-core/v5/common/errors/errorgen

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"hash/crc64"
	"math/rand"
	"sync"
	"time"

	core "github.com/v2fly/v2ray-core/v5"
	"github.com/v2fly/v2ray-core/v5/common"
	"github.com/v2fly/v2ray-core/v5/common/buf"
	"github.com/v2fly/v2ray-core/v5/common/net"
	"github.com/v2fly/v2ray-core/v5/common/platform"
	"github.com/v2fly/v2ray-core/v5/common/protocol"
	"github.com/v2fly/v2ray-core/v5/common/retry"
	"github.com/v2fly/v2ray-core/v5/common/serial"
	"github.com/v2fly/v2ray-core/v5/common/session"
	"github.com/v2fly/v2ray-core/v5/common/signal"
	"github.com/v2fly/v2ray-core/v5/common/task"
	"github.com/v2fly/v2ray-core/v5/features/policy"
	"github.com/v2fly/v2ray-core/v5/proxy"
	"github.com/v2fly/v2ray-core/v5/proxy/vmess"
	"github.com/v2fly/v2ray-core/v5/proxy/vmess/encoding"
	"github.com/v2fly/v2ray-core/v5/transport"
	"github.com/v2fly/v2ray-core/v5/transport/internet"
)

// [BDI] RTTStats maintains a moving average of the Round Trip Time
type RTTStats struct {
	sync.Mutex
	avg time.Duration
}

func (s *RTTStats) Update(val time.Duration) {
	s.Lock()
	defer s.Unlock()
	if s.avg == 0 {
		s.avg = val
	} else {
		// Exponential Moving Average with alpha=0.2
		s.avg = time.Duration(float64(s.avg)*0.8 + float64(val)*0.2)
	}
}

func (s *RTTStats) Get() time.Duration {
	s.Lock()
	defer s.Unlock()
	if s.avg == 0 {
		return 300 * time.Millisecond // Default pessimistic guess
	}
	return s.avg
}

// Handler is an outbound connection handler for VMess protocol.
type Handler struct {
	serverList    *protocol.ServerList
	serverPicker  protocol.ServerPicker
	policyManager policy.Manager
	rttStats      *RTTStats // [BDI] Estimator for Delta_Remote
}

// New creates a new VMess outbound handler.
func New(ctx context.Context, config *Config) (*Handler, error) {
	serverList := protocol.NewServerList()
	for _, rec := range config.Receiver {
		s, err := protocol.NewServerSpecFromPB(rec)
		if err != nil {
			return nil, newError("failed to parse server spec").Base(err)
		}
		serverList.AddServer(s)
	}

	v := core.MustFromContext(ctx)
	handler := &Handler{
		serverList:    serverList,
		serverPicker:  protocol.NewRoundRobinServerPicker(serverList),
		policyManager: v.GetFeature(policy.ManagerType()).(policy.Manager),
		rttStats:      &RTTStats{}, // [BDI] Initialize estimator
	}

	return handler, nil
}

// [BDI] sendDummyRequest establishes a separate connection to inject fake upstream traffic.
func (h *Handler) sendDummyRequest(ctx context.Context, dialer internet.Dialer, rec *protocol.ServerSpec) {
	// 1. Establish connection
	conn, err := dialer.Dial(ctx, rec.Destination())
	if err != nil {
		return // Best effort, ignore errors
	}
	defer conn.Close()

	// 2. Prepare Dummy Header
	user := rec.PickUser()
	account := user.Account.(*vmess.MemoryAccount)

	request := &protocol.RequestHeader{
		Version:  encoding.Version,
		User:     user,
		Command:  protocol.RequestCommandTCP, // Command doesn't matter much as IsDummy is true
		Address:  net.LocalHostIP,            // Fake address
		Port:     0,
		Option:   protocol.RequestOptionChunkStream,
		Security: account.Security,
		IsDummy:  true, // [BDI] CRITICAL: Flag this as a dummy packet
	}

	if request.Security == protocol.SecurityType_ZERO {
		request.Security = protocol.SecurityType_NONE
	}

	// 3. Initialize Crypto Session
	isAEAD := !aeadDisabled && len(account.AlterIDs) == 0
	hashkdf := hmac.New(sha256.New, []byte("VMessBF"))
	hashkdf.Write(account.ID.Bytes())
	behaviorSeed := crc64.Checksum(hashkdf.Sum(nil), crc64.MakeTable(crc64.ISO))
	session := encoding.NewClientSession(ctx, isAEAD, protocol.DefaultIDHash, int64(behaviorSeed))

	// 4. Write Data
	writer := buf.NewBufferedWriter(buf.NewWriter(conn))

	// Encode Header
	if err := session.EncodeRequestHeader(request, writer); err != nil {
		return
	}

	// Encode Random Body (Size randomization: 100 - 1500 bytes)
	bodyWriter, err := session.EncodeRequestBody(request, writer)
	if err != nil {
		return
	}

	dummySize := rand.Intn(1400) + 100
	dummyData := make([]byte, dummySize)
	rand.Read(dummyData)

	b := buf.New()
	b.Write(dummyData)
	bodyWriter.WriteMultiBuffer(buf.MultiBuffer{b})

	writer.Flush()
	// Connection closes on defer
}

// Process implements proxy.Outbound.Process().
func (h *Handler) Process(ctx context.Context, link *transport.Link, dialer internet.Dialer) error {
	var rec *protocol.ServerSpec
	var conn internet.Connection

	err := retry.ExponentialBackoff(5, 200).On(func() error {
		rec = h.serverPicker.PickServer()
		rawConn, err := dialer.Dial(ctx, rec.Destination())
		if err != nil {
			return err
		}
		conn = rawConn

		return nil
	})
	if err != nil {
		return newError("failed to find an available destination").Base(err).AtWarning()
	}
	defer conn.Close()

	outbound := session.OutboundFromContext(ctx)
	if outbound == nil || !outbound.Target.IsValid() {
		return newError("target not specified").AtError()
	}

	target := outbound.Target
	newError("tunneling request to ", target, " via ", rec.Destination().NetAddr()).WriteToLog(session.ExportIDToError(ctx))

	command := protocol.RequestCommandTCP
	if target.Network == net.Network_UDP {
		command = protocol.RequestCommandUDP
	}
	if target.Address.Family().IsDomain() && target.Address.Domain() == "v1.mux.cool" {
		command = protocol.RequestCommandMux
	}

	user := rec.PickUser()
	request := &protocol.RequestHeader{
		Version: encoding.Version,
		User:    user,
		Command: command,
		Address: target.Address,
		Port:    target.Port,
		Option:  protocol.RequestOptionChunkStream,
	}

	account := request.User.Account.(*vmess.MemoryAccount)
	request.Security = account.Security

	if request.Security == protocol.SecurityType_AES128_GCM || request.Security == protocol.SecurityType_NONE || request.Security == protocol.SecurityType_CHACHA20_POLY1305 {
		request.Option.Set(protocol.RequestOptionChunkMasking)
	}

	if shouldEnablePadding(request.Security) && request.Option.Has(protocol.RequestOptionChunkMasking) {
		request.Option.Set(protocol.RequestOptionGlobalPadding)
	}

	if request.Security == protocol.SecurityType_ZERO {
		request.Security = protocol.SecurityType_NONE
		request.Option.Clear(protocol.RequestOptionChunkStream)
		request.Option.Clear(protocol.RequestOptionChunkMasking)
	}

	if account.AuthenticatedLengthExperiment {
		request.Option.Set(protocol.RequestOptionAuthenticatedLength)
	}

	input := link.Reader
	output := link.Writer

	isAEAD := false
	if !aeadDisabled && len(account.AlterIDs) == 0 {
		isAEAD = true
	}

	hashkdf := hmac.New(sha256.New, []byte("VMessBF"))
	hashkdf.Write(account.ID.Bytes())

	behaviorSeed := crc64.Checksum(hashkdf.Sum(nil), crc64.MakeTable(crc64.ISO))

	session := encoding.NewClientSession(ctx, isAEAD, protocol.DefaultIDHash, int64(behaviorSeed))
	sessionPolicy := h.policyManager.ForLevel(request.User.Level)

	ctx, cancel := context.WithCancel(ctx)
	timer := signal.CancelAfterInactivity(ctx, cancel, sessionPolicy.Timeouts.ConnectionIdle)

	// [BDI] Start Time Measurement for Estimator
	requestStartTime := time.Now()

	// [BDI] Signal channel to cancel dummy packet if real response arrives
	headerReceived := make(chan struct{})

	requestDone := func() error {
		defer timer.SetTimeout(sessionPolicy.Timeouts.DownlinkOnly)

		writer := buf.NewBufferedWriter(buf.NewWriter(conn))
		if err := session.EncodeRequestHeader(request, writer); err != nil {
			return newError("failed to encode request").Base(err).AtWarning()
		}

		bodyWriter, err := session.EncodeRequestBody(request, writer)
		if err != nil {
			return newError("failed to start encoding").Base(err)
		}
		if err := buf.CopyOnceTimeout(input, bodyWriter, proxy.FirstPayloadTimeout); err != nil && err != buf.ErrNotTimeoutReader && err != buf.ErrReadTimeout {
			return newError("failed to write first payload").Base(err)
		}

		if err := writer.SetBuffered(false); err != nil {
			return err
		}

		// [BDI] Scheduler for Fake Uplink Request
		// Trigger logic: Wait for (EstimatedRTT * factor), then fire if no response yet.
		go func() {
			avgRTT := h.rttStats.Get()
			// Wait between 50% and 80% of the estimated RTT
			waitRatio := 0.5 + rand.Float64()*0.3
			waitTime := time.Duration(float64(avgRTT) * waitRatio)

			timer := time.NewTimer(waitTime)
			defer timer.Stop()

			select {
			case <-headerReceived:
				// Real response arrived, no need for dummy
				return
			case <-ctx.Done():
				// Connection closed
				return
			case <-timer.C:
				// Time is up, real response still missing. Inject Fake Request.
				h.sendDummyRequest(context.Background(), dialer, rec)
			}
		}()

		if err := buf.Copy(input, bodyWriter, buf.UpdateActivity(timer)); err != nil {
			return err
		}

		if request.Option.Has(protocol.RequestOptionChunkStream) && !account.NoTerminationSignal {
			if err := bodyWriter.WriteMultiBuffer(buf.MultiBuffer{}); err != nil {
				return err
			}
		}

		return nil
	}

	responseDone := func() error {
		defer timer.SetTimeout(sessionPolicy.Timeouts.UplinkOnly)

		reader := &buf.BufferedReader{Reader: buf.NewReader(conn)}
		header, err := session.DecodeResponseHeader(reader)
		if err != nil {
			return newError("failed to read header").Base(err)
		}

		// [BDI] Mark header as received to cancel/prevent dummy packet
		close(headerReceived)

		// [BDI] Update Estimator with actual RTT
		h.rttStats.Update(time.Since(requestStartTime))

		h.handleCommand(rec.Destination(), header.Command)

		bodyReader, err := session.DecodeResponseBody(request, reader)
		if err != nil {
			return newError("failed to start encoding response").Base(err)
		}
		return buf.Copy(bodyReader, output, buf.UpdateActivity(timer))
	}

	responseDonePost := task.OnSuccess(responseDone, task.Close(output))
	if err := task.Run(ctx, requestDone, responseDonePost); err != nil {
		return newError("connection ends").Base(err)
	}

	return nil
}

var (
	enablePadding = false
	aeadDisabled  = false
)

func shouldEnablePadding(s protocol.SecurityType) bool {
	return enablePadding || s == protocol.SecurityType_AES128_GCM || s == protocol.SecurityType_CHACHA20_POLY1305 || s == protocol.SecurityType_AUTO
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*Config))
	}))

	common.Must(common.RegisterConfig((*SimplifiedConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		simplifiedClient := config.(*SimplifiedConfig)
		fullClient := &Config{Receiver: []*protocol.ServerEndpoint{
			{
				Address: simplifiedClient.Address,
				Port:    simplifiedClient.Port,
				User: []*protocol.User{
					{
						Account: serial.ToTypedMessage(&vmess.Account{Id: simplifiedClient.Uuid}),
					},
				},
			},
		}}

		return common.CreateObject(ctx, fullClient)
	}))

	const defaultFlagValue = "NOT_DEFINED_AT_ALL"

	paddingValue := platform.NewEnvFlag("v2ray.vmess.padding").GetValue(func() string { return defaultFlagValue })
	if paddingValue != defaultFlagValue {
		enablePadding = true
	}

	isAeadDisabled := platform.NewEnvFlag("v2ray.vmess.aead.disabled").GetValue(func() string { return defaultFlagValue })
	if isAeadDisabled == "true" {
		aeadDisabled = true
	}
}
