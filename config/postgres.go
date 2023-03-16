package config

import "errors"

type Postgres struct {
	Path     string `mapstructure:"path"`
	DSN      string `mapstructure:"dsn"`
	DB       string `mapstructure:"db"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

func (p Postgres) GetDSN() string {
	if p.DSN != "" {
		return p.DSN
	}

	return ""
}

func (p Postgres) Validate() error {
	return errors.New("not implemented")
}
