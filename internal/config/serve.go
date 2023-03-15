package config

type Servers struct {
	Metrics *Server `mapstructure:"metrics"`
	Admin   *Server `mapstructure:"admin"`
	Public  *Server `mapstructure:"public"`
	Healthz *Server `mapstructure:"healthz"`
}

type Server struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	TLS  *TLS   `mapstructure:"tls"`
}

type TLS struct {
	CACert File `mapstructure:"ca-cert"`
	Cert   File `mapstructure:"cert"`
	Key    File `mapstructure:"key"`
}

type File struct {
	Path   string `mapstructure:"path"`
	Base64 string `mapstructure:"base64"`
}
