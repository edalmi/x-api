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
	"github.com/prometheus/client_golang/prometheus"
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

	options := &internal.Options{
		Cache:   cache,
		Pubsub:  pubsub,
		Queue:   queue,
		Logger:  logger,
		Metrics: prometheus.NewRegistry(),
	}

	users := handler.NewUser(options)
	groups := handler.NewGroup(options)

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
		cache:  cache,
		logger: logger,
		pubsub: pubsub,
		queue:  queue,

		public:  srvPublic,
		admin:   srvAdmin,
		metrics: srvMetrics,
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

	public  *http.Server
	admin   *http.Server
	metrics *http.Server
	healthz *http.Server
}

func (s *Server) Start() error {
	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGINT)

	go func() {
		if err := s.admin.ListenAndServe(); err != nil {
			s.logger.Error(err)
			return
		}
	}()

	go func() {
		if err := s.public.ListenAndServe(); err != nil {
			s.logger.Error(err)
			return
		}
	}()

	go func() {
		if err := s.metrics.ListenAndServe(); err != nil {
			s.logger.Error(err)
			return
		}
	}()

	go func() {
		if err := s.healthz.ListenAndServe(); err != nil {
			s.logger.Error(err)
			return
		}
	}()

	defer s.Close()

	<-c

	if err := s.metrics.Shutdown(context.Background()); err != nil {
		s.logger.Warn(err)
	}

	if err := s.public.Shutdown(context.Background()); err != nil {
		s.logger.Warn(err)
	}

	if err := s.admin.Shutdown(context.Background()); err != nil {
		s.logger.Warn(err)
	}

	return nil
}

func (s *Server) Close() error {
	if err := s.admin.Close(); err != nil {
		s.logger.Error(err)
	}

	if err := s.public.Close(); err != nil {
		s.logger.Error(err)
	}

	if err := s.metrics.Close(); err != nil {
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
