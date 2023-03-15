package config

import (
	"github.com/redis/go-redis/v9"
)

type Cache struct {
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
