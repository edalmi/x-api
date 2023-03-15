package server

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/edalmi/x-api"
	"github.com/edalmi/x-api/config"
	"github.com/edalmi/x-api/handler"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func New(cfg *config.Config) (*Server, error) {
	cache, err := setupCache(cfg.Cache)
	if err != nil {
		return nil, err
	}

	logger, err := setupLogger(cfg.Logger)
	if err != nil {
		return nil, err
	}

	pubsub, err := setupPubsub(cfg.Pubsub)
	if err != nil {
		return nil, err
	}

	queue, err := setupQueue(cfg.Queue)
	if err != nil {
		return nil, err
	}

	options := &xapi.Options{
		Cache:   cache,
		Pubsub:  pubsub,
		Queue:   queue,
		Logger:  logger,
		Metrics: prometheus.NewRegistry(),
	}

	srv := &Server{
		cfg:     cfg,
		options: options,
	}

	if err := srv.setupAdminServer(); err != nil {
		return nil, err
	}

	if err := srv.setupPublicServer(); err != nil {
		return nil, err
	}

	if err := srv.setupHealthzServer(); err != nil {
		return nil, err
	}

	if err := srv.setupMetrcisServer(); err != nil {
		return nil, err
	}

	return srv, nil
}

func (s Server) setupHealthzServer() error {
	handler := handler.NewHealthz(s.options)

	router := chi.NewRouter()
	router.Mount("/healthz", handler.Routes())

	srv, err := setupHTTPServer(s.cfg.Serve.Healthz, router)
	if err != nil {
		return err
	}

	s.healthz = srv

	return nil
}

func (s Server) setupMetrcisServer() error {
	handler := promhttp.HandlerFor(
		s.options.Metrics.(*prometheus.Registry),
		promhttp.HandlerOpts{
			Registry: s.options.Metrics,
		},
	)

	srv, err := setupHTTPServer(s.cfg.Serve.Metrics, handler)
	if err != nil {
		return err
	}

	s.metrics = srv

	return nil
}

func (s Server) setupPublicServer() error {
	var (
		users  = handler.NewUser(s.options)
		groups = handler.NewGroup(s.options)
	)

	router := chi.NewRouter()
	router.Mount("/users", users.Routes())
	router.Mount("/groups", groups.Routes())

	srv, err := setupHTTPServer(s.cfg.Serve.Public, router)
	if err != nil {
		return err
	}

	s.public = srv

	return nil
}

func (s Server) setupAdminServer() error {
	var (
		users  = handler.NewUser(s.options)
		groups = handler.NewGroup(s.options)
	)

	router := chi.NewRouter()
	router.Mount("/users", users.Routes())
	router.Mount("/groups", groups.Routes())

	srv, err := setupHTTPServer(s.cfg.Serve.Admin, router)
	if err != nil {
		return err
	}

	s.public = srv

	return nil
}

type publicRouter interface {
	RegisterRoutes(r *chi.Router)
}

type adminRouter interface {
	RegisterAdminRoutes(r *chi.Router)
}

type Server struct {
	cfg     *config.Config
	options *xapi.Options
	public  *HTTPServer
	admin   *HTTPServer
	metrics *HTTPServer
	healthz *HTTPServer
}

func (s *Server) Start() error {
	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGINT)

	go func() {
		if s.public.TLS {
			if err := s.public.ListenAndServeTLS(s.public.TLSCert, s.public.TLSKey); err != nil {
				s.options.Logger.Error(err)
				return
			}
		} else {
			if err := s.public.ListenAndServe(); err != nil {
				s.options.Logger.Error(err)
				return
			}
		}
	}()

	go func() {
		if s.public.TLS {
			if err := s.admin.ListenAndServeTLS(s.admin.TLSCert, s.admin.TLSKey); err != nil {
				s.options.Logger.Error(err)
				return
			}
		} else {
			if err := s.admin.ListenAndServe(); err != nil {
				s.options.Logger.Error(err)
				return
			}
		}
	}()

	go func() {
		if s.metrics.TLS {
			if err := s.metrics.ListenAndServeTLS(s.metrics.TLSCert, s.metrics.TLSKey); err != nil {
				s.options.Logger.Error(err)
				return
			}
		} else {
			if err := s.metrics.ListenAndServe(); err != nil {
				s.options.Logger.Error(err)
				return
			}
		}
	}()

	go func() {
		if s.healthz.TLS {
			if err := s.healthz.ListenAndServeTLS(s.healthz.TLSCert, s.healthz.TLSKey); err != nil {
				s.options.Logger.Error(err)
				return
			}
		} else {
			if err := s.healthz.ListenAndServe(); err != nil {
				s.options.Logger.Error(err)
				return
			}
		}
	}()

	defer s.Close()

	<-c

	if err := s.metrics.Shutdown(context.Background()); err != nil {
		s.options.Logger.Warn(err)
	}

	if err := s.public.Shutdown(context.Background()); err != nil {
		s.options.Logger.Warn(err)
	}

	if err := s.admin.Shutdown(context.Background()); err != nil {
		s.options.Logger.Warn(err)
	}

	if err := s.healthz.Shutdown(context.Background()); err != nil {
		s.options.Logger.Warn(err)
	}

	return nil
}

func (s *Server) Close() error {
	if err := s.closeOptions(); err != nil {
		return err
	}

	if err := s.admin.Close(); err != nil {
		s.options.Logger.Error(err)
	}

	if err := s.public.Close(); err != nil {
		s.options.Logger.Error(err)
	}

	if err := s.metrics.Close(); err != nil {
		s.options.Logger.Error(err)
	}

	if err := s.healthz.Close(); err != nil {
		s.options.Logger.Error(err)
	}

	return nil
}

func (s *Server) closeOptions() error {
	if c, ok := s.options.Cache.(io.Closer); ok {
		c.Close()
	}

	if c, ok := s.options.Pubsub.(io.Closer); ok {
		c.Close()
	}

	if c, ok := s.options.Queue.(io.Closer); ok {
		c.Close()
	}

	if c, ok := s.options.Logger.(io.Closer); ok {
		c.Close()
	}

	return nil
}
