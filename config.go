package ginja

import "github.com/imdario/mergo"

type Config struct {
	Namespace  string
	Version    string
	MountStats bool `json:"-"`
	StatsURL   string
	Debug      bool
}

var (
	defaultConfig = Config{
		Namespace:  "api",
		Version:    "1",
		MountStats: false,
		StatsURL:   "_stats",
		Debug:      true,
	}
)

func (c *Config) buildUrl() string {
	return "/" + c.Namespace + "/v" + c.Version
}

func (c *Config) ApplyDefaults() {
	// config := defaultConfig

	setDebug := !c.Debug
	var shouldDebug bool
	if setDebug {
		shouldDebug = c.Debug
	}

	mergo.MergeWithOverwrite(c, defaultConfig)

	if setDebug {
		c.Debug = shouldDebug
	}
}
