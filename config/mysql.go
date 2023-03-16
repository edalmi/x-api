package config

import "errors"

type MySQL struct {
	Path     string `mapstructure:"path"`
	DSN      string `mapstructure:"dsn"`
	DB       string `mapstructure:"db"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

func (p MySQL) GetDSN() string {
	if p.DSN != "" {
		return p.DSN
	}

	return ""
}

func (p MySQL) Validate() error {
	return errors.New("not implemented")
}
