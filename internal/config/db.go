package config

type DB struct {
	Path     string `mapstructure:"path"`
	DSN      string `mapstructure:"dsn"`
	DB       string `mapstructure:"db"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
}
