package server

import (
	"net/http"

	"github.com/edalmi/x-api/config"
)

func setupHTTPServer(cfg *config.Server, h http.Handler) (*http.Server, error) {
	return &http.Server{
		Handler: h,
	}, nil
}
