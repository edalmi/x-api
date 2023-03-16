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
		app:    cfg.App,
		config: cfg,
		logger: stdlog.New(log.Default()),
	}

	if err := srv.setupLogger(); err != nil {
		return nil, err
	}

	if err := srv.setupDB(); err != nil {
		return nil, err
	}

	if err := srv.setupCache(); err != nil {
		return nil, err
	}

	if err := srv.setupMetrics(); err != nil {
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
	s.logger.Info("setup metrics")

	s.prometheus = prometheus.NewRegistry()

	return nil
}

func (s *Server) setupLogger() error {
	logger, err := setupLogger(s.config.Mode, s.config.Logger)
	if err != nil {
		return err
	}

	s.logger = logger
	return nil
}

func (s *Server) setupDB() error {
	s.logger.Info("setup database")

	db, err := setupDB(s.config.DB)
	if err != nil {
		return err
	}

	s.db = db
	return nil
}

func (s *Server) setupCache() error {
	s.logger.Info("setup database")

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
	app        string
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

func (s Server) App() string {
	return s.app
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
		s.logger.Infof("Starting public server at %v", s.public.Addr)
		s.startServer(s.public)
	}()

	go func() {
		s.logger.Infof("Starting admin at %v", s.admin.Addr)
		s.startServer(s.admin)
	}()

	go func() {
		s.logger.Infof("Starting metrics server at %v", s.metrics.Addr)
		s.startServer(s.metrics)
	}()

	go func() {
		s.logger.Infof("Starting healthz server at %v", s.healthz.Addr)
		s.startServer(s.healthz)
	}()

	defer func() {
		s.logger.Info("Tearing down public server")
		if err := s.teardownServer(s.public); err != nil {
			s.logger.Error(err)
		}

		s.logger.Info("Tearing down admin server")
		if err := s.teardownServer(s.admin); err != nil {
			s.logger.Error(err)
		}

		s.logger.Info("Tearing down metrics server")
		if err := s.teardownServer(s.metrics); err != nil {
			s.logger.Error(err)
		}

		s.logger.Info("Tearing down healthz server")
		if err := s.teardownServer(s.healthz); err != nil {
			s.logger.Error(err)
		}

		s.logger.Info("Tearing down cache provider")
		if err := s.teardownCache(); err != nil {
			s.logger.Error(err)
		}

		s.logger.Info("Tearing down queue provider")
		if err := s.teardownQueue(); err != nil {
			s.logger.Error(err)
		}

		s.logger.Info("Tearing down pubsub provider")
		if err := s.teardownPubsub(); err != nil {
			s.logger.Error(err)
		}

		s.logger.Info("Tearing down queue provider")
		if err := s.teardownQueue(); err != nil {
			s.logger.Error(err)
		}

		s.logger.Info("Tearing down database provider")
		if err := s.teardownDB(); err != nil {
			s.logger.Error(err)
		}

		s.logger.Info("Tearing down logger provider")
		if err := s.teardownLogger(); err != nil {
			s.logger.Error(err)
		}
	}()

	<-c

	return nil
}

func (s Server) startServer(srv *HTTPServer) error {
	s.logger.Infof("Starting healthz server at %v", s.healthz.Addr)

	if srv.TLS {
		if err := s.healthz.ListenAndServeTLS(srv.TLSCert, srv.TLSKey); err != nil {
			if err == http.ErrServerClosed {
				return nil
			}

			return err
		}
	} else {
		if err := srv.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				return nil
			}

			return err
		}
	}

	return nil
}

func (s *Server) teardownServer(srv *HTTPServer) error {
	if err := srv.Shutdown(context.Background()); err != nil {
		if err == http.ErrServerClosed {
			return nil
		}

		return err
	}

	if err := srv.Close(); err != nil {
		if err == http.ErrServerClosed {
			return nil
		}

		return err
	}

	return nil
}

func (s *Server) teardownCache() error {
	if c, ok := s.cache.(io.Closer); ok {
		return c.Close()
	}

	return nil
}

func (s *Server) teardownPubsub() error {
	if c, ok := s.pubsub.(io.Closer); ok {
		return c.Close()
	}

	return nil
}

func (s *Server) teardownQueue() error {
	if c, ok := s.queue.(io.Closer); ok {
		return c.Close()
	}

	return nil
}

func (s *Server) teardownLogger() error {
	if c, ok := s.logger.(io.Closer); ok {
		return c.Close()
	}

	return nil
}

func (s *Server) teardownDB() error {
	if c, ok := s.logger.(io.Closer); ok {
		return c.Close()
	}

	return nil
}
