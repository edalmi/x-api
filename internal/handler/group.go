package handler

import (
	"net/http"

	"github.com/edalmi/x-api/internal"
	"github.com/go-chi/chi/v5"
)

func NewGroup(_ *internal.Options) *Group {
	return &Group{}
}

type Group struct {
	opts internal.Options
}

func (u Group) CreateGroup(rw http.ResponseWriter, r *http.Request) {}

func (u Group) ListGroups(rw http.ResponseWriter, r *http.Request) {}

func (u Group) DeleteGroup(rw http.ResponseWriter, r *http.Request) {}

func (u Group) GetGroup(rw http.ResponseWriter, r *http.Request) {}

func (u Group) PublicRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", u.ListGroups)
	r.Get("/:id", u.GetGroup)
	r.Post("/", u.CreateGroup)
	r.Delete("/", u.DeleteGroup)

	return r
}
