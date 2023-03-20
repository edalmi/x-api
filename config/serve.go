package config

import (
	"errors"
)

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

func (s Server) Validate() error {
	return errors.New("not implemented")
}

type TLS struct {
	Cert string `mapstructure:"cert"`
	Key  string `mapstructure:"key"`
}

func (t TLS) Validate() error {
	return errors.New("not implemented")
}
