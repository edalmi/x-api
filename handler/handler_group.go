package handler

import (
	"net/http"

	"github.com/edalmi/x-api/internal"
	"github.com/go-chi/chi/v5"
)

func NewGroup(_ Options) *GroupHandler {
	return &GroupHandler{}
}

type GroupHandler struct {
	opts internal.Options
}

func (u GroupHandler) CreateGroup(rw http.ResponseWriter, r *http.Request) {}

func (u GroupHandler) ListGroups(rw http.ResponseWriter, r *http.Request) {}

func (u GroupHandler) DeleteGroup(rw http.ResponseWriter, r *http.Request) {}

func (u GroupHandler) GetGroup(rw http.ResponseWriter, r *http.Request) {}

func (u GroupHandler) Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", u.ListGroups)
	r.Get("/:id", u.GetGroup)
	r.Post("/", u.CreateGroup)
	r.Delete("/", u.DeleteGroup)

	return r
}
