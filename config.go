package ginja

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
	return c.Namespace + "/v" + c.Version
}
