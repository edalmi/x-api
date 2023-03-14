package server

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/edalmi/x-api/internal"
	"github.com/edalmi/x-api/internal/config"
	"github.com/edalmi/x-api/internal/handler"
)

func New(cfg *config.Config) (*Server, error) {
	cache, err := configureCache(cfg.Cache)
	if err != nil {
		return nil, err
	}

	logger, err := configureLogger(cfg.Logger)
	if err != nil {
		return nil, err
	}

	pubsub, err := configurePubsub(cfg.Pubsub)
	if err != nil {
		return nil, err
	}

	queue, err := configureQueue(cfg.Queue)
	if err != nil {
		return nil, err
	}

	router := &Server{
		cache:  cache,
		logger: logger,
		pubsub: pubsub,
		queue:  queue,

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

	logger internal.Logger
	cache  internal.Cache
	pubsub internal.Pubsub
	queue  internal.Queue

	Users  *handler.User
	Groups *handler.Group
}

func (s *Server) Start() error {
	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGINT)

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

	<-c

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
