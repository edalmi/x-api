package config

import "errors"

type SQLite struct {
	Path     string `mapstructure:"path"`
	DSN      string `mapstructure:"dsn"`
	InMemory string `mapstructure:"in-memory"`
}

func (p SQLite) GetDSN() string {
	if p.DSN != "" {
		return p.DSN
	}

	return ""
}

func (p SQLite) Validate() error {
	return errors.New("not implemented")
}
