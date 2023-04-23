package apiserver

type Config struct {
	BindAddr   string `toml:"bind_addr"`
	SessionKey string `toml:"session_key"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
	}
}
