package handler

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/edalmi/x-api/database"
	"github.com/edalmi/x-api/database/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
)

func NewUserHandler(opts HandlerOptions) *UserHandler {
	return &UserHandler{
		UserRepository: postgres.UserRepo{
			DB: opts.DB(),
		},
		UserMetrics: newUserMetrics(opts.Metrics()),
		Options:     opts,
	}
}

type UserHandler struct {
	UserRepository userRepo
	UserMetrics    UserMetrics
	Options        HandlerOptions
}

func (u *UserHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	_, err := u.UserRepository.CreateUser(context.Background(), database.NewUser{})
	if err != nil {
		return
	}

	u.UserMetrics.IncrementUsersCreated()
}

func (u UserHandler) ListUsers(rw http.ResponseWriter, r *http.Request) {
	log := log.New(os.Stdout, "users", 0)
	log.SetPrefix("users")
	log.Println(r.URL.Path)

	_, err := u.UserRepository.ListUsers(context.Background())
	if err != nil {
		return
	}

	u.UserMetrics.IncrementUsersCreated()
}

func (u UserHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	log := log.New(os.Stdout, "users", 0)
	log.SetPrefix("users")
	log.Println(r.URL.Path)

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

type userRepo interface {
	CreateUser(ctx context.Context, in database.NewUser) (*database.User, error)
	ListUsers(ctx context.Context) ([]database.User, error)
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
