package handler

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
)

func NewUserHandler(opts HandlerOptions) *UserHandler {
	return &UserHandler{
		metrics: newUserMetrics(opts.Metrics()),
		opts:    opts,
	}
}

type UserHandler struct {
	metrics *userMetrics
	opts    HandlerOptions
}

func (u *UserHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	u.metrics.incrementCreatedUsers()
}

func (u UserHandler) ListUsers(rw http.ResponseWriter, r *http.Request) {
	log := log.New(os.Stdout, "users", 0)
	log.SetPrefix("users")
	log.Println(r.URL.Path)

	u.metrics.incrementCreatedUsers()
}

func (u UserHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	log := log.New(os.Stdout, "users", 0)
	log.SetPrefix("users")
	log.Println(r.URL.Path)

	u.metrics.incrementDeletedUsers()
}

func (u UserHandler) GetUser(rw http.ResponseWriter, r *http.Request) {}

func (u UserHandler) UpdateUser(rw http.ResponseWriter, r *http.Request) {}

func (u UserHandler) Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", u.ListUsers)
	r.Get("/{id}", u.GetUser)
	r.Post("/", u.CreateUser)
	r.Delete("/{id}", u.DeleteUser)
	r.Put("/{id}", u.UpdateUser)

	return r
}

type userMetrics struct {
	mu           sync.Mutex
	createdUsers prometheus.Counter
	deletedUsers prometheus.Counter
}

func (u *userMetrics) incrementCreatedUsers() {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.createdUsers.Inc()
}

func (u *userMetrics) incrementDeletedUsers() {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.deletedUsers.Inc()
}

func newUserMetrics(reg prometheus.Registerer) *userMetrics {
	m := &userMetrics{
		createdUsers: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "x",
			Name:      "users_created",
			Help:      "Number of created users",
		}),
		deletedUsers: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "x",
			Name:      "users_deleted",
			Help:      "Number of deleted users",
		}),
	}

	reg.MustRegister(m.createdUsers, m.deletedUsers)

	return m
}
