package inbound

// GetDefaultValue returns default settings of DefaultConfig.
func (c *Config) GetDefaultValue() *DefaultConfig {
	if c.GetDefault() == nil {
		return &DefaultConfig{
			Level:   0,
			AlterId: 0,
		}
	}
	return c.Default
}

// [BDI] GetBDISettings returns the BDI configuration or safe defaults.
func (c *Config) GetBDISettings() *BDIConfig {
	if c.GetBdiSettings() == nil {
		return &BDIConfig{
			Enabled:    false,
			PaddingMin: 100,
			PaddingMax: 1000,
			JitterMin:  5,
			JitterMax:  50,
		}
	}
	return c.BdiSettings
}
