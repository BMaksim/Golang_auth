package api

type Config struct {
	Addr string `json:"addr"`
}

func NewConfig() *Config {
	return &Config{
		Addr: ":8080",
	}
}
