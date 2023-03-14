package server

import (
	"io"
	"log"
	"net/http"

	"github.com/edalmi/x-api/internal"
	"github.com/edalmi/x-api/internal/config"
	"github.com/edalmi/x-api/internal/handler"
)

func New(cfg *config.Config) (*Server, error) {
	cache, err := configureCache(cfg.Cache)
	if err != nil {
		return nil, err
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
	publicRouter *http.Server
	adminRouter  *http.Server
	cache        internal.Cache

	Users  *handler.User
	Groups *handler.Group
}

func (s *Server) Listen() error {
	go func() {
		if err := s.publicRouter.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	go func() {
		if err := s.adminRouter.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	defer s.Close()

	return nil
}

func (s *Server) Close() error {
	s.publicRouter.Close()
	s.adminRouter.Close()

	if c, ok := s.cache.(io.Closer); ok {
		c.Close()
	}

	return nil
}
