package outbound

// [BDI] GetBDISettings returns the BDI configuration or safe defaults.
func (c *Config) GetBDISettings() *BDIConfig {
	if c.GetBdiSettings() == nil {
		return &BDIConfig{
			Enabled:    false,
			PacketRate: 2, // Default to 2 dummy packets per second if enabled
		}
	}
	return c.BdiSettings
}
