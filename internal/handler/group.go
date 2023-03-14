package handler

import (
	"net/http"

	"github.com/edalmi/x-api/internal"
)

type Group struct {
	Cache internal.Cache
}

func (u Group) CreateGroup(rw http.ResponseWriter, r *http.Request) {}

func (u Group) ListGroups(rw http.ResponseWriter, r *http.Request) {}

func (u Group) DeleteGroup(rw http.ResponseWriter, r *http.Request) {}

func (u Group) GetGroup(rw http.ResponseWriter, r *http.Request) {}

func (u Group) SetPublicRoute() {}
