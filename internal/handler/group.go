package handler

import (
	"net/http"

	"github.com/edalmi/x-api"
	"github.com/go-chi/chi/v5"
)

func NewGroup(_ *xapi.Options) *Group {
	return &Group{}
}

type Group struct {
	opts xapi.Options
}

func (u Group) CreateGroup(rw http.ResponseWriter, r *http.Request) {}

func (u Group) ListGroups(rw http.ResponseWriter, r *http.Request) {}

func (u Group) DeleteGroup(rw http.ResponseWriter, r *http.Request) {}

func (u Group) GetGroup(rw http.ResponseWriter, r *http.Request) {}

func (u Group) Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", u.ListGroups)
	r.Get("/:id", u.GetGroup)
	r.Post("/", u.CreateGroup)
	r.Delete("/", u.DeleteGroup)

	return r
}
