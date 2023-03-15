package config

type Serve struct {
	Metrics *ServeItem `mapstructure:"metrics"`
	Admin   *ServeItem `mapstructure:"admin"`
	Public  *ServeItem `mapstructure:"public"`
	Healthz *ServeItem `mapstructure:"healthz"`
}

type ServeItem struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	TLS  *TLS   `mapstructure:"tls"`
}

type TLS struct {
	CACert File `mapstructure:"cacert"`
	Cert   File `mapstructure:"cert"`
	Key    File `mapstructure:"key"`
}

type File struct {
	Path   string `mapstructure:"path"`
	Base64 string `mapstructure:"base64"`
}
