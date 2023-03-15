package server

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/edalmi/x-api/internal"
	"github.com/edalmi/x-api/internal/config"
	"github.com/edalmi/x-api/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	users := &handler.User{
		Cache: cache,
	}

	groups := &handler.Group{
		Cache: cache,
	}

	apiv1 := chi.NewRouter()
	apiv1.Mount("/users", users.PublicRoutes())
	apiv1.Mount("/groups", groups.PublicRoutes())

	srvPublic, err := configureHTTP(cfg.Serve.Public, apiv1)
	if err != nil {
		return nil, err
	}

	adminV1 := chi.NewRouter()
	srvAdmin, err := configureHTTP(cfg.Serve.Admin, adminV1)
	if err != nil {
		return nil, err
	}

	metrics := http.NewServeMux()
	metrics.HandleFunc("/", promhttp.Handler())
	srvMetrics, err := configureHTTP(cfg.Serve.Metrics, metrics)
	if err != nil {
		return nil, err
	}

	srv := &Server{
		cache:      cache,
		logger:     logger,
		pubsub:     pubsub,
		queue:      queue,
		srvPublic:  srvPublic,
		srvAdmin:   srvAdmin,
		srvMetrics: srvMetrics,
	}

	return srv, nil
}

type publicRouter interface {
	RegisterRoutes(r *chi.Router)
}

type adminRouter interface {
	RegisterAdminRoutes(r *chi.Router)
}

type Server struct {
	logger internal.Logger
	cache  internal.Cache
	pubsub internal.Pubsub
	queue  internal.Queue

	srvPublic  *http.Server
	srvAdmin   *http.Server
	srvMetrics *http.Server
}

func (s *Server) Start() error {
	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGINT)

	go func() {
		if err := s.srvAdmin.ListenAndServe(); err != nil {
			s.logger.Error(err)
			return
		}
	}()

	go func() {
		if err := s.srvPublic.ListenAndServe(); err != nil {
			s.logger.Error(err)
			return
		}
	}()

	go func() {
		if err := s.srvMetrics.ListenAndServe(); err != nil {
			s.logger.Error(err)
			return
		}
	}()

	defer s.Close()

	<-c

	if err := s.srvMetrics.Shutdown(context.Background()); err != nil {
		s.logger.Warn(err)
	}

	if err := s.srvPublic.Shutdown(context.Background()); err != nil {
		s.logger.Warn(err)
	}

	if err := s.srvAdmin.Shutdown(context.Background()); err != nil {
		s.logger.Warn(err)
	}

	return nil
}

func (s *Server) Close() error {
	if err := s.srvAdmin.Close(); err != nil {
		s.logger.Error(err)
	}

	if err := s.srvPublic.Close(); err != nil {
		s.logger.Error(err)
	}

	if err := s.srvMetrics.Close(); err != nil {
		s.logger.Error(err)
	}

	if c, ok := s.cache.(io.Closer); ok {
		c.Close()
	}

	if c, ok := s.logger.(io.Closer); ok {
		c.Close()
	}

	if c, ok := s.queue.(io.Closer); ok {
		c.Close()
	}

	if c, ok := s.pubsub.(io.Closer); ok {
		c.Close()
	}

	return nil
}
