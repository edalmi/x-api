package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewGroupHandler(_ HandlerOpts) *GroupHandler {
	return &GroupHandler{}
}

type GroupHandler struct {
	opts HandlerOpts
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

type GroupCreate struct{}

type Group struct{}

type GroupService interface {
	CreateUser(ctx context.Context, g GroupCreate) (*Group, error)
}
