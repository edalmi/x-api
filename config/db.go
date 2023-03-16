package config

import "errors"

type DB struct {
	Postgres *Postgres `mapstructure:"postgres"`
	MySQL    *MySQL    `mapstructure:"mysql"`
	SQLite   *SQLite   `mapstructure:"sqlite"`
	MariaDB  *MariaDB  `mapstructure:"mariadb"`
}

func (p DB) Validate() error {
	parts := validator{
		"postgres": p.Postgres,
		"mysql":    p.MySQL,
		"sqlite":   p.SQLite,
		"mariadb":  p.MariaDB,
	}

	return parts.Validate()
}

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

type MariaDB struct {
	Path     string `mapstructure:"path"`
	DSN      string `mapstructure:"dsn"`
	DB       string `mapstructure:"db"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

func (p MariaDB) GetDSN() string {
	if p.DSN != "" {
		return p.DSN
	}

	return ""
}

func (p MariaDB) Validate() error {
	return errors.New("not implemented")
}
