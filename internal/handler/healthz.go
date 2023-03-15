package handler

import (
	"net/http"

	"github.com/edalmi/x-api/internal"
	"github.com/go-chi/chi/v5"
)

type Healthz struct {
	Cache internal.Cache
}

func (u Healthz) Check(rw http.ResponseWriter, r *http.Request) {}

func (u Healthz) Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", u.Check)

	return r
}
