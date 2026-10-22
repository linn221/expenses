package config

type cfg struct {
	api_port       string
	base_directory string
	dev_mode       bool
}

// singleton
var Cfg cfg

func (c *cfg) API_PORT() string {
	return c.api_port
}

func (c *cfg) BASE_DIR() string {
	return c.base_directory
}

func (c *cfg) DEV_MODE() bool {
	return c.dev_mode
}
