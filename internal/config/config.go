package config

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func New(v *viper.Viper) (*Config, error) {
	cfg := DefaultConfig()
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func DefaultConfig() Config {
	return Config{
		Admin: &Listen{
			Port: 11234,
		},
		Public: &Listen{
			Port: 11235,
		},
	}
}

type Listen struct {
	Port int `mapstructure:"port"`
}

type Config struct {
	Admin  *Listen
	Public *Listen
	Cache  *Cache `mapstructure:"cache"`
	TLS    *TLS   `mapstructure:"tls"`
}

type Cache struct {
	Provider  string     `mapstructure:"provider"`
	Redis     *Redis     `mapstructure:"redis"`
	Memcached *Memcached `mapstructure:"memcached"`
}

type Redis struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func (r Redis) Config() (*redis.Options, error) {
	return &redis.Options{
		Addr:     r.Address,
		Password: r.Password,
		DB:       r.DB,
	}, nil
}

type Memcached struct {
	Addresses []string `mapstructure:"addresses"`
}

type TLS struct {
	CACert []string `mapstructure:"ca-cert"`
	Cert   []string `mapstructure:"cert"`
	Key    []string `mapstructure:"key"`
}
