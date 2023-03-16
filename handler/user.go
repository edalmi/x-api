package handler

import (
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
)

var reqs int64

func NewUserHandler(opts HandlerOptions) *UserHandler {
	return &UserHandler{
		UserMetrics: newUserMetrics(opts.App(), opts.Metrics()),
		Options:     opts,
	}
}

type UserHandler struct {
	UserMetrics UserMetrics
	Options     HandlerOptions
}

func (u *UserHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	u.Options.Logger().Info(r.URL.Path)

	u.UserMetrics.IncrementUsersCreated()
}

func (u UserHandler) ListUsers(rw http.ResponseWriter, r *http.Request) {
	u.Options.Logger().Info(r.URL.Path)

	time.Sleep(time.Duration(reqs) * time.Second)

	atomic.AddInt64(&reqs, 100)

	u.UserMetrics.IncrementUsersCreated()
}

func (u UserHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	u.Options.Logger().Info(r.URL.Path)

	time.Sleep(3 * time.Second)

	u.UserMetrics.IncrementUsersDeleted()
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

type UserMetrics interface {
	IncrementUsersCreated()
	IncrementUsersDeleted()
}

type userMetrics struct {
	mu           sync.Mutex
	createdUsers prometheus.Counter
	deletedUsers prometheus.Counter
}

func (u *userMetrics) IncrementUsersCreated() {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.createdUsers.Inc()
}

func (u *userMetrics) IncrementUsersDeleted() {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.deletedUsers.Inc()
}

func newUserMetrics(app string, reg prometheus.Registerer) *userMetrics {
	m := &userMetrics{
		createdUsers: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: app,
			Name:      "users_created",
			Help:      "Number of created users",
		}),
		deletedUsers: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: app,
			Name:      "users_deleted",
			Help:      "Number of deleted users",
		}),
	}

	reg.MustRegister(m.createdUsers, m.deletedUsers)

	return m
}
