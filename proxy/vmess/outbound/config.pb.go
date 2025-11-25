package outbound

import (
	net "github.com/v2fly/v2ray-core/v5/common/net"
	protocol "github.com/v2fly/v2ray-core/v5/common/protocol"
	_ "github.com/v2fly/v2ray-core/v5/common/protoext"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// [BDI] Traffic Shaping Configuration for Client
type BDIConfig struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Enabled       bool                   `protobuf:"varint,1,opt,name=enabled,proto3" json:"enabled,omitempty"`
	PacketRate    int32                  `protobuf:"varint,2,opt,name=packet_rate,json=packetRate,proto3" json:"packet_rate,omitempty"` // Packets per second for fake requests (approx)
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BDIConfig) Reset() {
	*x = BDIConfig{}
	mi := &file_proxy_vmess_outbound_config_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BDIConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BDIConfig) ProtoMessage() {}

func (x *BDIConfig) ProtoReflect() protoreflect.Message {
	mi := &file_proxy_vmess_outbound_config_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BDIConfig.ProtoReflect.Descriptor instead.
func (*BDIConfig) Descriptor() ([]byte, []int) {
	return file_proxy_vmess_outbound_config_proto_rawDescGZIP(), []int{0}
}

func (x *BDIConfig) GetEnabled() bool {
	if x != nil {
		return x.Enabled
	}
	return false
}

func (x *BDIConfig) GetPacketRate() int32 {
	if x != nil {
		return x.PacketRate
	}
	return 0
}

type Config struct {
	state    protoimpl.MessageState     `protogen:"open.v1"`
	Receiver []*protocol.ServerEndpoint `protobuf:"bytes,1,rep,name=Receiver,proto3" json:"Receiver,omitempty"`
	// [BDI] Add BDI config to outbound
	BdiSettings   *BDIConfig `protobuf:"bytes,2,opt,name=bdi_settings,json=bdiSettings,proto3" json:"bdi_settings,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Config) Reset() {
	*x = Config{}
	mi := &file_proxy_vmess_outbound_config_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Config) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Config) ProtoMessage() {}

func (x *Config) ProtoReflect() protoreflect.Message {
	mi := &file_proxy_vmess_outbound_config_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Config.ProtoReflect.Descriptor instead.
func (*Config) Descriptor() ([]byte, []int) {
	return file_proxy_vmess_outbound_config_proto_rawDescGZIP(), []int{1}
}

func (x *Config) GetReceiver() []*protocol.ServerEndpoint {
	if x != nil {
		return x.Receiver
	}
	return nil
}

func (x *Config) GetBdiSettings() *BDIConfig {
	if x != nil {
		return x.BdiSettings
	}
	return nil
}

type SimplifiedConfig struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Address       *net.IPOrDomain        `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Port          uint32                 `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
	Uuid          string                 `protobuf:"bytes,3,opt,name=uuid,proto3" json:"uuid,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SimplifiedConfig) Reset() {
	*x = SimplifiedConfig{}
	mi := &file_proxy_vmess_outbound_config_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SimplifiedConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SimplifiedConfig) ProtoMessage() {}

func (x *SimplifiedConfig) ProtoReflect() protoreflect.Message {
	mi := &file_proxy_vmess_outbound_config_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SimplifiedConfig.ProtoReflect.Descriptor instead.
func (*SimplifiedConfig) Descriptor() ([]byte, []int) {
	return file_proxy_vmess_outbound_config_proto_rawDescGZIP(), []int{2}
}

func (x *SimplifiedConfig) GetAddress() *net.IPOrDomain {
	if x != nil {
		return x.Address
	}
	return nil
}

func (x *SimplifiedConfig) GetPort() uint32 {
	if x != nil {
		return x.Port
	}
	return 0
}

func (x *SimplifiedConfig) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

var File_proxy_vmess_outbound_config_proto protoreflect.FileDescriptor

const file_proxy_vmess_outbound_config_proto_rawDesc = "" +
	"\n" +
	"!proxy/vmess/outbound/config.proto\x12\x1fv2ray.core.proxy.vmess.outbound\x1a!common/protocol/server_spec.proto\x1a\x18common/net/address.proto\x1a common/protoext/extensions.proto\"F\n" +
	"\tBDIConfig\x12\x18\n" +
	"\aenabled\x18\x01 \x01(\bR\aenabled\x12\x1f\n" +
	"\vpacket_rate\x18\x02 \x01(\x05R\n" +
	"packetRate\"\x9f\x01\n" +
	"\x06Config\x12F\n" +
	"\bReceiver\x18\x01 \x03(\v2*.v2ray.core.common.protocol.ServerEndpointR\bReceiver\x12M\n" +
	"\fbdi_settings\x18\x02 \x01(\v2*.v2ray.core.proxy.vmess.outbound.BDIConfigR\vbdiSettings\"\x92\x01\n" +
	"\x10SimplifiedConfig\x12;\n" +
	"\aaddress\x18\x01 \x01(\v2!.v2ray.core.common.net.IPOrDomainR\aaddress\x12\x12\n" +
	"\x04port\x18\x02 \x01(\rR\x04port\x12\x12\n" +
	"\x04uuid\x18\x03 \x01(\tR\x04uuid:\x19\x82\xb5\x18\x15\n" +
	"\boutbound\x12\x05vmess\x90\xff)\x01B~\n" +
	"#com.v2ray.core.proxy.vmess.outboundP\x01Z3github.com/v2fly/v2ray-core/v5/proxy/vmess/outbound\xaa\x02\x1fV2Ray.Core.Proxy.Vmess.Outboundb\x06proto3"

var (
	file_proxy_vmess_outbound_config_proto_rawDescOnce sync.Once
	file_proxy_vmess_outbound_config_proto_rawDescData []byte
)

func file_proxy_vmess_outbound_config_proto_rawDescGZIP() []byte {
	file_proxy_vmess_outbound_config_proto_rawDescOnce.Do(func() {
		file_proxy_vmess_outbound_config_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proxy_vmess_outbound_config_proto_rawDesc), len(file_proxy_vmess_outbound_config_proto_rawDesc)))
	})
	return file_proxy_vmess_outbound_config_proto_rawDescData
}

var file_proxy_vmess_outbound_config_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_proxy_vmess_outbound_config_proto_goTypes = []any{
	(*BDIConfig)(nil),               // 0: v2ray.core.proxy.vmess.outbound.BDIConfig
	(*Config)(nil),                  // 1: v2ray.core.proxy.vmess.outbound.Config
	(*SimplifiedConfig)(nil),        // 2: v2ray.core.proxy.vmess.outbound.SimplifiedConfig
	(*protocol.ServerEndpoint)(nil), // 3: v2ray.core.common.protocol.ServerEndpoint
	(*net.IPOrDomain)(nil),          // 4: v2ray.core.common.net.IPOrDomain
}
var file_proxy_vmess_outbound_config_proto_depIdxs = []int32{
	3, // 0: v2ray.core.proxy.vmess.outbound.Config.Receiver:type_name -> v2ray.core.common.protocol.ServerEndpoint
	0, // 1: v2ray.core.proxy.vmess.outbound.Config.bdi_settings:type_name -> v2ray.core.proxy.vmess.outbound.BDIConfig
	4, // 2: v2ray.core.proxy.vmess.outbound.SimplifiedConfig.address:type_name -> v2ray.core.common.net.IPOrDomain
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_proxy_vmess_outbound_config_proto_init() }
func file_proxy_vmess_outbound_config_proto_init() {
	if File_proxy_vmess_outbound_config_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proxy_vmess_outbound_config_proto_rawDesc), len(file_proxy_vmess_outbound_config_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proxy_vmess_outbound_config_proto_goTypes,
		DependencyIndexes: file_proxy_vmess_outbound_config_proto_depIdxs,
		MessageInfos:      file_proxy_vmess_outbound_config_proto_msgTypes,
	}.Build()
	File_proxy_vmess_outbound_config_proto = out.File
	file_proxy_vmess_outbound_config_proto_goTypes = nil
	file_proxy_vmess_outbound_config_proto_depIdxs = nil
}
