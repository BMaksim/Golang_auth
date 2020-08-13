package api

type Config struct {
	Addr         string `json:"addr"`
	DatabaseURL  string `json:"database_url"`
	DatabaseName string `json:"database_name"`
}

func NewConfig() *Config {
	return &Config{
		Addr: ":8080",
	}
}
