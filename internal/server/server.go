package server

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/edalmi/x-api/internal"
	"github.com/edalmi/x-api/internal/config"
	"github.com/edalmi/x-api/internal/handler"
	stdlog "github.com/edalmi/x-api/internal/logging/log"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func New(cfg *config.Config) (*Server, error) {
	log.Println("setup cache")
	cache, err := setupCache(cfg.Cache)
	if err != nil {
		return nil, err
	}

	/*	log.Println("setup logger")
		logger, err := setupLogger(cfg.Logger)
		if err != nil {
			return nil, err
		}

		log.Println("setup pubsub")
		pubsub, err := setupPubsub(cfg.Pubsub)
		if err != nil {
			return nil, err
		}

		log.Println("setup queue")
		queue, err := setupQueue(cfg.Queue)
		if err != nil {
			return nil, err
		}*/

	options := &internal.Options{
		Cache: cache,
		// Pubsub:  pubsub,
		// Queue:   queue,
		Logger: &stdlog.Logger{
			Logger: log.Default(),
		},
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

func (s *Server) setupHealthzServer() error {
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

func (s *Server) setupMetrcisServer() error {
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

func (s *Server) setupPublicServer() error {
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

func (s *Server) setupAdminServer() error {
	router := http.NewServeMux()
	srv, err := setupHTTPServer(s.cfg.Serve.Admin, router)
	if err != nil {
		return err
	}

	s.admin = srv

	return nil
}

type Server struct {
	cfg     *config.Config
	options *internal.Options
	public  *HTTPServer
	admin   *HTTPServer
	metrics *HTTPServer
	healthz *HTTPServer
}

func (s *Server) Start() error {
	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGINT)

	s.options.Logger.Info("PID:", os.Getpid())

	go func() {
		log.Printf("Starting public server at %v", s.public.Addr)

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
		log.Printf("Starting admin at %v", s.admin.Addr)

		if s.admin.TLS {
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
		log.Printf("Starting metrics server at %v", s.metrics.Addr)

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
		log.Printf("Starting healthz server at %v", s.healthz.Addr)

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

	defer func() {
		s.options.Logger.Info("Shutting down metrics server")
		if err := s.metrics.Shutdown(context.Background()); err != nil {
			s.options.Logger.Warn(err)
		}

		s.options.Logger.Info("Shutting down public server")
		if err := s.public.Shutdown(context.Background()); err != nil {
			s.options.Logger.Warn("X", err)
		}

		s.options.Logger.Info("Shutting down admin server")
		if err := s.admin.Shutdown(context.Background()); err != nil {
			s.options.Logger.Warn(err)
		}

		s.options.Logger.Info("Shutting down healthz server")
		if err := s.healthz.Shutdown(context.Background()); err != nil {
			s.options.Logger.Warn(err)
		}

		s.cleanUp()
	}()

	<-c

	return nil
}

func (s *Server) cleanUp() error {
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
		s.options.Logger.Info("Freeing cache resources")
		c.Close()
	}

	if c, ok := s.options.Pubsub.(io.Closer); ok {
		s.options.Logger.Info("Freeing pubsub resources")
		c.Close()
	}

	if c, ok := s.options.Queue.(io.Closer); ok {
		s.options.Logger.Info("Freeing queue resources")
		c.Close()
	}

	if c, ok := s.options.Logger.(io.Closer); ok {
		s.options.Logger.Info("Freeing logger resources")
		c.Close()
	}

	return nil
}
