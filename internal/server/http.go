package server

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/edalmi/x-api/internal/config"
)

func setupHTTPServer(cfg *config.Server, handler http.Handler) (*HTTPServer, error) {
	srv := &HTTPServer{
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

		srv.TLS = true
		srv.TLSCert = tlsCfg.Cert
		srv.TLSKey = tlsCfg.Key
	}

	return srv, nil
}

type HTTPServer struct {
	*http.Server
	TLS     bool
	TLSCert string
	TLSKey  string
}
