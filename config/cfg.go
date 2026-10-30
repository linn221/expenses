package config

import "time"

type cfg struct {
	api_port       string
	base_directory string
	dev_mode       bool
	lockTTL        time.Duration
}

// singleton
var Cfg cfg

func (c *cfg) API_PORT() string {
	if c.api_port == "" {
		return "8080"
	}
	return c.api_port
}

func (c *cfg) BASE_DIR() string {
	return c.base_directory
}

func (c *cfg) DEV_MODE() bool {
	return c.dev_mode
}

func (c *cfg) LockTTL() time.Duration {
	return time.Minute * 2
}

func (c *cfg) PageSize() int {
	return 20
}

func (c *cfg) SessionTTL() time.Duration {
	return time.Hour * 24 * 7
}

func (c *cfg) RateLimitWindow() time.Duration {
	return time.Minute * 2
}

func (c *cfg) RateLimit() int64 {
	return 50
}

func (c *cfg) SlowQueryThreshold() time.Duration {
	return time.Second
}
