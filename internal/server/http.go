package server

import (
	"errors"
	"net/http"

	"github.com/edalmi/x-api/internal/config"
)

func configureHTTP(cfg *config.ServerItem, h http.Handler) (*http.Server, error) {
	return &http.Server{
		Handler: h,
	}, nil
}
