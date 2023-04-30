package proxyserver

type Config struct {
	BindAddrProxy string `toml:"bind_addr_proxy"`
	BindAddrApi   string `toml:"bind_addr_api"`
	BindAddrHtml  string `toml:"bind_addr_html"`
	SessionKey    string `toml:"session_key"`
}

func NewConfig() *Config {
	return &Config{
		BindAddrProxy: ":8080",
		BindAddrApi:   ":8081",
		BindAddrHtml:  ":8082",
	}
}
