package server

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/edalmi/x-api/config"
)

func setupHTTPServer(cfg *config.Server, handler http.Handler) (*httpServer, error) {
	srv := &httpServer{
		Server: &http.Server{
			Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Handler:      handler,
			ReadTimeout:  time.Duration(cfg.ReadTimeout),
			WriteTimeout: time.Duration(cfg.WriteTimeout),
		},
	}

	if tlsCfg := cfg.TLS; tlsCfg != nil {
		if tlsCfg.Cert == "" {
			return nil, errors.New("cert error")
		}

		if tlsCfg.Key == "" {
			return nil, errors.New("key error")
		}

		srv.tls = true
		srv.tlsCert = tlsCfg.Cert
		srv.tlsKey = tlsCfg.Key
	}

	return srv, nil
}

type httpServer struct {
	*http.Server
	tls     bool
	tlsCert string
	tlsKey  string
}
