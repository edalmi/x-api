package router

import (
	"errors"

	"github.com/edalmi/x-api/internal"
	"github.com/edalmi/x-api/internal/config"
	"github.com/edalmi/x-api/internal/handler"
	memcachedprovider "github.com/edalmi/x-api/internal/memcached"
	redisprovider "github.com/edalmi/x-api/internal/redis"
	"github.com/redis/go-redis/v9"
)

func NewRouter(cfg *config.Config) (*Router, error) {
	if cfg.Cache == nil {
		// Use in memory cache
	}

	var cache internal.Cache
	if cfg := cfg.Cache; cfg.Provider == "redis" && cfg.Redis != nil {
		redisCfg, err := cfg.Redis.Config()
		if err != nil {
			return nil, err
		}

		cache = redisprovider.NewCache(redis.NewClient(redisCfg))
	}

	if cfg := cfg.Cache; cfg.Provider == "memcached" && cfg.Memcached != nil {
		addr := cfg.Memcached.Addresses
		if len(addr) == 0 {
			return nil, errors.New("no addresses")
		}

		mcCache, err := memcachedprovider.New(addr)
		if err != nil {
			return nil, err
		}

		cache = mcCache
	}

	router := &Router{
		Users: &handler.User{
			Cache: cache,
		},
		Groups: &handler.Group{
			Cache: cache,
		},
	}

	return router, nil
}

type Router struct {
	Users  *handler.User
	Groups *handler.Group
}
