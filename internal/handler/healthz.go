package handler

import (
	"net/http"

	"github.com/edalmi/x-api/internal"
	"github.com/go-chi/chi/v5"
)

func NewHealthz(_ *internal.Options) *Healthz {
	return &Healthz{}
}

type Healthz struct {
	opts internal.Options
}

func (u Healthz) Check(rw http.ResponseWriter, r *http.Request) {}

func (u Healthz) Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", u.Check)

	return r
}
