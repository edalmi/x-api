package server

import (
	"errors"

	"github.com/edalmi/x-api/caching"
	memcachedprovider "github.com/edalmi/x-api/caching/memcached"
	redisprovider "github.com/edalmi/x-api/caching/redis"
	"github.com/edalmi/x-api/config"
	"github.com/redis/go-redis/v9"
)

func setupCache(cfg *config.Cache) (caching.Cache, error) {
	if cfg == nil {
		return nil, errors.New("cache is empty")
	}

	if cfg.Redis != nil {
		redisCfg, err := cfg.Redis.Config()
		if err != nil {
			return nil, err
		}

		return redisprovider.NewCache(redis.NewClient(redisCfg)), nil
	}

	if cfg.Memcached != nil {
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
