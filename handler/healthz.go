package handler

import (
	"net/http"

	"github.com/edalmi/x-api"
	"github.com/go-chi/chi/v5"
)

func NewHealthz(_ *xapi.Options) *Healthz {
	return &Healthz{}
}

type Healthz struct {
	opts xapi.Options
}

func (u Healthz) Live(rw http.ResponseWriter, r *http.Request) {}

func (u Healthz) Ready(rw http.ResponseWriter, r *http.Request) {}

func (u Healthz) Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/live", u.Live)
	r.Get("/ready", u.Ready)

	return r
}
