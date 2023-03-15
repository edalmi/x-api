package handler

import (
	"net/http"
	"sync"

	"github.com/edalmi/x-api/internal"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
)

func NewUser(opts *internal.Options) *User {
	return &User{
		metrics: NewUserMetrics(opts.Metrics),
		opts:    opts,
	}
}

type User struct {
	metrics *UserMetrics
	opts    *internal.Options
}

func (u *User) CreateUser(rw http.ResponseWriter, r *http.Request) {
	u.metrics.IncTotalUsers()
}

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

type UserMetrics struct {
	mu    sync.Mutex
	Users prometheus.Counter
}

func (u *UserMetrics) IncTotalUsers() {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.Users.Inc()
}

func NewUserMetrics(reg prometheus.Registerer) *UserMetrics {
	m := &UserMetrics{
		Users: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "x",
			Name:      "users_created",
			Help:      "Number of created users",
		}),
	}

	reg.MustRegister(m.Users)

	return m
}
