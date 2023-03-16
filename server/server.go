package server

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/edalmi/x-api/caching"
	"github.com/edalmi/x-api/config"
	"github.com/edalmi/x-api/database"
	"github.com/edalmi/x-api/handler"
	"github.com/edalmi/x-api/logging"
	stdlog "github.com/edalmi/x-api/logging/log"
	"github.com/edalmi/x-api/pubsub"
	"github.com/edalmi/x-api/queue"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func New(cfg *config.Config) (*Server, error) {
	srv := &Server{
		config: cfg,
	}

	if err := srv.setupDB(); err != nil {
		return nil, err
	}

	if err := srv.setupCache(); err != nil {
		return nil, err
	}

	if err := srv.setupLogger(); err != nil {
		return nil, err
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

func (s *Server) setupMetrics() error {
	log.Println("setup metrics")

	s.logger = &stdlog.Logger{
		Logger: log.Default(),
	}

	return nil
}

func (s *Server) setupLogger() error {
	log.Println("setup metrics")
	s.prometheus = prometheus.NewRegistry()

	return nil
}

func (s *Server) setupDB() error {
	log.Println("setup database")
	db, err := setupDB(s.config.DB)
	if err != nil {
		return err
	}

	s.db = db
	return nil
}

func (s *Server) setupCache() error {
	log.Println("setup database")

	cache, err := setupCache(s.config.Cache)
	if err != nil {
		return err
	}

	s.cache = cache
	return nil
}

func (s *Server) setupHealthzServer() error {
	handler := handler.NewHealthz(s)

	router := chi.NewRouter()
	router.Mount("/healthz", handler.Routes())

	srv, err := setupHTTPServer(s.config.Serve.Healthz, router)
	if err != nil {
		return err
	}

	s.healthz = srv

	return nil
}

func (s *Server) setupMetrcisServer() error {
	handler := promhttp.HandlerFor(
		s.prometheus.(*prometheus.Registry),
		promhttp.HandlerOpts{
			Registry: s.prometheus,
		},
	)

	srv, err := setupHTTPServer(s.config.Serve.Metrics, handler)
	if err != nil {
		return err
	}

	s.metrics = srv

	return nil
}

func (s *Server) setupPublicServer() error {
	var (
		users  = handler.NewUserHandler(s)
		groups = handler.NewGroupHandler(s)
	)

	router := chi.NewRouter()
	router.Mount("/users", users.Routes())
	router.Mount("/groups", groups.Routes())

	srv, err := setupHTTPServer(s.config.Serve.Public, router)
	if err != nil {
		return err
	}

	s.public = srv

	return nil
}

func (s *Server) setupAdminServer() error {
	router := http.NewServeMux()
	srv, err := setupHTTPServer(s.config.Serve.Admin, router)
	if err != nil {
		return err
	}

	s.admin = srv

	return nil
}

type Server struct {
	config     *config.Config
	db         *database.DB
	cache      caching.Cache
	logger     logging.Logger
	pubsub     pubsub.Pubsub
	queue      queue.Queue
	prometheus prometheus.Registerer

	public  *HTTPServer
	admin   *HTTPServer
	metrics *HTTPServer
	healthz *HTTPServer
}

func (s Server) Config() *config.Config {
	return s.config
}

func (s Server) Logger() logging.Logger {
	return s.logger
}

func (s Server) Cache() caching.Cache {
	return s.cache
}

func (s Server) Queue() queue.Queue {
	return s.queue
}

func (s Server) Pubsub() pubsub.Pubsub {
	return s.pubsub
}

func (s Server) DB() *database.DB {
	return s.db
}

func (s Server) Metrics() prometheus.Registerer {
	return s.prometheus
}

func (s *Server) Start() error {
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	s.logger.Info("PID:", os.Getpid())

	go func() {
		log.Printf("Starting public server at %v", s.public.Addr)

		if s.public.TLS {
			if err := s.public.ListenAndServeTLS(s.public.TLSCert, s.public.TLSKey); err != nil {
				s.logger.Error(err)
				return
			}
		} else {
			if err := s.public.ListenAndServe(); err != nil {
				s.logger.Error(err)
				return
			}
		}
	}()

	go func() {
		log.Printf("Starting admin at %v", s.admin.Addr)

		if s.admin.TLS {
			if err := s.admin.ListenAndServeTLS(s.admin.TLSCert, s.admin.TLSKey); err != nil {
				s.logger.Error(err)
				return
			}
		} else {
			if err := s.admin.ListenAndServe(); err != nil {
				s.logger.Error(err)
				return
			}
		}
	}()

	go func() {
		log.Printf("Starting metrics server at %v", s.metrics.Addr)

		if s.metrics.TLS {
			if err := s.metrics.ListenAndServeTLS(s.metrics.TLSCert, s.metrics.TLSKey); err != nil {
				s.logger.Error(err)
				return
			}
		} else {
			if err := s.metrics.ListenAndServe(); err != nil {
				s.logger.Error(err)
				return
			}
		}
	}()

	go func() {
		log.Printf("Starting healthz server at %v", s.healthz.Addr)

		if s.healthz.TLS {
			if err := s.healthz.ListenAndServeTLS(s.healthz.TLSCert, s.healthz.TLSKey); err != nil {
				s.logger.Error(err)
				return
			}
		} else {
			if err := s.healthz.ListenAndServe(); err != nil {
				s.logger.Error(err)
				return
			}
		}
	}()

	defer func() {
		s.logger.Info("Shutting down metrics server")
		if err := s.metrics.Shutdown(context.Background()); err != nil {
			s.logger.Warn(err)
		}

		s.logger.Info("Shutting down public server")
		if err := s.public.Shutdown(context.Background()); err != nil {
			s.logger.Warn("X", err)
		}

		s.logger.Info("Shutting down admin server")
		if err := s.admin.Shutdown(context.Background()); err != nil {
			s.logger.Warn(err)
		}

		s.logger.Info("Shutting down healthz server")
		if err := s.healthz.Shutdown(context.Background()); err != nil {
			s.logger.Warn(err)
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
		s.logger.Error(err)
	}

	if err := s.public.Close(); err != nil {
		s.logger.Error(err)
	}

	if err := s.metrics.Close(); err != nil {
		s.logger.Error(err)
	}

	if err := s.healthz.Close(); err != nil {
		s.logger.Error(err)
	}

	return nil
}

func (s *Server) closeOptions() error {
	if c, ok := s.cache.(io.Closer); ok {
		s.logger.Info("Freeing cache resources")
		c.Close()
	}

	if c, ok := s.pubsub.(io.Closer); ok {
		s.logger.Info("Freeing pubsub resources")
		c.Close()
	}

	if c, ok := s.queue.(io.Closer); ok {
		s.logger.Info("Freeing queue resources")
		c.Close()
	}

	if c, ok := s.logger.(io.Closer); ok {
		s.logger.Info("Freeing logger resources")
		c.Close()
	}

	return nil
}
