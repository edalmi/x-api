package server

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/edalmi/x-api/caching"
	"github.com/edalmi/x-api/config"
	"github.com/edalmi/x-api/database"
	"github.com/edalmi/x-api/handler"
	"github.com/edalmi/x-api/logging"
	stdlog "github.com/edalmi/x-api/logging/log"
	"github.com/edalmi/x-api/pubsub"
	"github.com/edalmi/x-api/queue"
	"github.com/go-chi/chi/v5"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/sync/errgroup"
)

func New(cfg *config.Config) (*Server, error) {
	srv := &Server{
		id:     cfg.App,
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

	if err := srv.setupPrometheus(); err != nil {
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

	if err := srv.setupOtel(); err != nil {
		return nil, err
	}

	return srv, nil
}

func (s *Server) setupPrometheus() error {
	s.logger.Info("setting up metrics provider")
	s.prometheus = prom.NewRegistry()

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
	s.logger.Info("setting up database")

	db, err := setupDB(s.config.DB)
	if err != nil {
		return err
	}

	s.db = db
	return nil
}

func (s *Server) setupOtel() error {
	t, err := setupOtel()
	if err != nil {
		return err
	}

	s.tracing = t
	otel.SetTracerProvider(s.tracing)
	return nil
}

func (s *Server) setupCache() error {
	s.logger.Info("setting up cache provider")

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

	s.healthzServer = srv

	return nil
}

func (s *Server) setupMetrcisServer() error {
	handler := promhttp.HandlerFor(
		s.prometheus.(*prom.Registry),
		promhttp.HandlerOpts{
			Registry: s.prometheus,
		},
	)

	srv, err := setupHTTPServer(s.config.Serve.Metrics, handler)
	if err != nil {
		return err
	}

	s.metricsServer = srv

	return nil
}

func (s *Server) setupPublicServer() error {
	var (
		usersHandler  = handler.NewUserHandler(s)
		groupsHandler = handler.NewGroupHandler(s)
	)

	router := chi.NewRouter()

	router.Mount("/users", usersHandler.Routes())
	router.Mount("/groups", groupsHandler.Routes())

	srv, err := setupHTTPServer(s.config.Serve.Public, router)
	if err != nil {
		return err
	}

	s.publicServer = srv

	return nil
}

func (s *Server) setupAdminServer() error {
	router := http.NewServeMux()
	srv, err := setupHTTPServer(s.config.Serve.Admin, router)
	if err != nil {
		return err
	}

	s.adminServer = srv

	return nil
}

type Server struct {
	tracing    *sdktrace.TracerProvider
	id         string
	config     *config.Config
	db         *database.DB
	cache      caching.Cache
	logger     logging.Logger
	pubsub     pubsub.Pubsub
	queue      queue.Queue
	prometheus prom.Registerer
	httpServers
}

type httpServers struct {
	publicServer  *httpServer
	adminServer   *httpServer
	metricsServer *httpServer
	healthzServer *httpServer
}

func (s Server) ID() string {
	return s.id
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

func (s Server) Prometheus() prom.Registerer {
	return s.prometheus
}

func (srv *Server) Start(ctx context.Context) error {
	sig := make(chan os.Signal, 1)

	signal.Notify(sig, os.Interrupt)

	srv.logger.Infof("PID: %d", os.Getpid())
	srv.logger.Infof("OS: %v/%v", runtime.GOOS, runtime.GOARCH)

	g := new(errgroup.Group)

	g.Go(func() error {
		srv.logger.Infof("Starting public server at %v", srv.publicServer.Addr)
		return srv.startHTTPServer(srv.publicServer)
	})

	g.Go(func() error {
		srv.logger.Infof("Starting admin at %v", srv.adminServer.Addr)
		return srv.startHTTPServer(srv.adminServer)
	})

	g.Go(func() error {
		srv.logger.Infof("Starting metrics server at %v", srv.metricsServer.Addr)
		return srv.startHTTPServer(srv.metricsServer)
	})

	g.Go(func() error {
		srv.logger.Infof("Starting healthz server at %v", srv.healthzServer.Addr)
		return srv.startHTTPServer(srv.healthzServer)
	})

	go func() {
		if err := g.Wait(); err != nil {
			srv.logger.Error(err)
		}
	}()

	defer func() {
		srv.logger.Info("Tearing down public server")
		if err := srv.shutdownHTTPServer(
			srv.publicServer,
			time.Duration(srv.config.Serve.Public.ShutdownTimeout),
		); err != nil {
			srv.logger.Error(err)
		}

		srv.logger.Info("Tearing down admin server")
		if err := srv.shutdownHTTPServer(
			srv.adminServer,
			time.Duration(srv.config.Serve.Admin.ShutdownTimeout),
		); err != nil {
			srv.logger.Error(err)
		}

		srv.logger.Info("Tearing down metrics server")
		if err := srv.shutdownHTTPServer(
			srv.metricsServer,
			time.Duration(srv.config.Serve.Metrics.ShutdownTimeout),
		); err != nil {
			srv.logger.Error(err)
		}

		srv.logger.Info("Tearing down healthz server")
		if err := srv.shutdownHTTPServer(
			srv.healthzServer,
			time.Duration(srv.config.Serve.Healthz.ShutdownTimeout),
		); err != nil {
			srv.logger.Error(err)
		}

		srv.logger.Info("Tearing down cache provider")
		if err := release(srv.cache); err != nil {
			srv.logger.Error(err)
		}

		srv.logger.Info("Tearing down pubsub provider")
		if err := release(srv.pubsub); err != nil {
			srv.logger.Error(err)
		}

		srv.logger.Info("Tearing down queue provider")
		if err := release(srv.queue); err != nil {
			srv.logger.Error(err)
		}

		srv.logger.Info("Tearing down database provider")
		if err := release(srv.db); err != nil {
			srv.logger.Error(err)
		}

		srv.logger.Info("Tearing down logger provider")
		if err := release(srv.logger); err != nil {
			srv.logger.Error(err)
		}
	}()

	<-sig
	srv.logger.Info("Shutting down servers")

	return nil
}

func (s Server) startHTTPServer(srv *httpServer) error {
	if srv.useTLS {
		err := srv.ListenAndServeTLS(srv.tlsCert, srv.tlsKey)
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		return err
	}

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return err
}

func (s *Server) shutdownHTTPServer(srv *httpServer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return err
	}

	return nil
}

func (s *Server) releaseOtel() error {
	if s.db != nil {
		return s.tracing.Shutdown(context.Background())
	}

	return nil
}

func release(resource interface{}) error {
	if c, ok := resource.(io.Closer); ok {
		return c.Close()
	}

	return nil
}
