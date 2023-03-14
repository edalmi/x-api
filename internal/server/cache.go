package server

import (
	"errors"

	"github.com/edalmi/x-api/internal"
	"github.com/edalmi/x-api/internal/config"
	memcachedprovider "github.com/edalmi/x-api/internal/memcached"
	redisprovider "github.com/edalmi/x-api/internal/redis"
	"github.com/redis/go-redis/v9"
)

const (
	cacheRedis     = "redis"
	cacheMemcached = "memcached"
	cacheInMemory  = "inmemory"
)

func configureCache(cfg *config.Cache) (internal.Cache, error) {
	if cfg == nil {
		cfg.Provider = cacheInMemory
	}

	if cfg.Provider == cacheRedis {
		if cfg.Redis == nil {
			return nil, errors.New("error")
		}

		redisCfg, err := cfg.Redis.Config()
		if err != nil {
			return nil, err
		}

		return redisprovider.NewCache(redis.NewClient(redisCfg)), nil
	}

	if cfg.Provider == cacheMemcached {
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
