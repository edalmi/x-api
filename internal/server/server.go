package server

import (
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/edalmi/x-api/internal"
	"github.com/edalmi/x-api/internal/config"
	"github.com/edalmi/x-api/internal/handler"
	memcachedprovider "github.com/edalmi/x-api/internal/memcached"
	redisprovider "github.com/edalmi/x-api/internal/redis"
	"github.com/redis/go-redis/v9"
)

func New(cfg *config.Config) (*Server, error) {
	if cfg.Cache == nil {
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

	router := &Server{
		cache: cache,

		Users: &handler.User{
			Cache: cache,
		},
		Groups: &handler.Group{
			Cache: cache,
		},
	}

	return router, nil
}

type Server struct {
	adminSrv  *http.Server
	publicSrv *http.Server
	cache     internal.Cache

	Users  *handler.User
	Groups *handler.Group
}

func (s *Server) Listen() error {
	go func() {
		if err := s.adminSrv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	go func() {
		if err := s.publicSrv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	defer s.Close()

	return nil
}

func (s *Server) Close() error {
	s.adminSrv.Close()
	s.publicSrv.Close()

	if c, ok := s.cache.(io.Closer); ok {
		c.Close()
	}

	return nil
}
