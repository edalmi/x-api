package server

import (
	"errors"
	"net/http"

	"github.com/edalmi/x-api/internal/config"
)

func configureHTTP(cfg *config.Server, h http.Handler) (*http.Server, error) {
	return &http.Server{
		Handler: h,
	}, nil
}
