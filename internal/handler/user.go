package handler

import (
	"net/http"

	"github.com/edalmi/x-api/internal"
)

type User struct {
	Cache internal.Cache
}

func (u User) CreateUser(rw http.ResponseWriter, r *http.Request) {}

func (u User) ListUsers(rw http.ResponseWriter, r *http.Request) {}

func (u User) DeleteUser(rw http.ResponseWriter, r *http.Request) {}

func (u User) GetUser(rw http.ResponseWriter, r *http.Request) {}

func (u User) SetPublicRoute() {}
