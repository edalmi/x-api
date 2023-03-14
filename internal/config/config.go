package config

import "github.com/redis/go-redis/v9"

type Config struct {
	Cache *Cache `mapstructure:"cache"`
	TLS   *TLS   `mapstructure:"tls"`
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
	}
}

type Memcached struct {
	Addresses []string `mapstructure:"addresses"`
}

type TLS struct {
	CACert []string `mapstructure:"ca_cert"`
	Cert   []string `mapstructure:"cert"`
	Key    []string `mapstructure:"key"`
}
