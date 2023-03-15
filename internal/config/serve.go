package config

type Servers struct {
	Metrics *Server `mapstructure:"metrics"`
	Admin   *Server `mapstructure:"admin"`
	Public  *Server `mapstructure:"public"`
	Healthz *Server `mapstructure:"healthz"`
}

type Server struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	TLS          *TLS   `mapstructure:"tls"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

type TLS struct {
	Cert string `mapstructure:"cert"`
	Key  string `mapstructure:"key"`
}
