package handler

import (
	"net/http"

	"github.com/edalmi/x-api/internal"
	"github.com/go-chi/chi/v5"
)

type User struct {
	Cache internal.Cache
}

func (u User) CreateUser(rw http.ResponseWriter, r *http.Request) {}

func (u User) ListUsers(rw http.ResponseWriter, r *http.Request) {}

func (u User) DeleteUser(rw http.ResponseWriter, r *http.Request) {}

func (u User) GetUser(rw http.ResponseWriter, r *http.Request) {}

func (u User) UpdateUser(rw http.ResponseWriter, r *http.Request) {}

func (u User) PublicRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", u.ListUsers)
	r.Get("/:id", u.GetUser)
	r.Post("/", u.CreateUser)
	r.Delete("/", u.DeleteUser)
	r.Put("/", u.UpdateUser)

	return r
}
