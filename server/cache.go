package server

import (
	"errors"

	"github.com/edalmi/x-api"
	"github.com/edalmi/x-api/xapi.config"
	memcachedprovider "github.com/edalmi/x-api/xapi.memcached"
	redisprovider "github.com/edalmi/x-api/xapi.redis"
	"github.com/redis/go-redis/v9"
)

func configureCache(cfg *config.Cache) (xapi.Cache, error) {
	if cfg == nil {
	}

	if cfg.Redis != nil {
		if cfg.Redis == nil {
			return nil, errors.New("error")
		}

		redisCfg, err := cfg.Redis.Config()
		if err != nil {
			return nil, err
		}

		return redisprovider.NewCache(redis.NewClient(redisCfg)), nil
	}

	if cfg.Memcached != nil {
		if cfg.Memcached == nil {
			return nil, errors.New("error")
		}

		addr := cfg.Memcached.Addresses
		if len(addr) == 0 {
			return nil, errors.New("no addresses")
		}

		memcached, err := memcachedprovider.New(addr)
		if err != nil {
			return nil, err
		}

		return memcached, nil
	}

	return nil, errors.New("error")
}
