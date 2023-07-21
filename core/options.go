package core

type Option func(*Config)

// Config 配置
func WithConfig(cfg *Config) Option {
	return func(c *Config) {
		if cfg.ModuleName != "" {
			c.ModuleName = cfg.ModuleName
		}
		if cfg.LoadNum != 0 {
			c.LoadNum = cfg.LoadNum
		}
		if cfg.LoadThreshold != 0 {
			c.LoadThreshold = cfg.LoadThreshold
		}
		if cfg.IncrFirstID != 0 {
			c.IncrFirstID = cfg.IncrFirstID
		}
	}
}

// WithLoadNum 设置每批生成数量
func WithLoadNum(num int32) Option {
	return func(cfg *Config) {
		cfg.LoadNum = num
	}
}

// WithLoadThreshold 设置加载阀值
func WithLoadThreshold(threshold int32) Option {
	return func(cfg *Config) {
		cfg.LoadThreshold = threshold
	}
}
