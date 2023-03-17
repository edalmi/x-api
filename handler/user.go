package handler

import (
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var reqs int64

func NewUserHandler(opts HandlerOptions) *UserHandler {
	return &UserHandler{
		UserMetrics: newUserMetrics(opts.ID(), opts.Metrics()),
		Options:     opts,
	}
}

type UserHandler struct {
	UserMetrics UserMetrics
	Options     HandlerOptions
}

func (u *UserHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	_, span := otel.Tracer(u.Options.ID()).Start(r.Context(), "users.CreateUser")
	defer span.End()

	u.Options.Logger().Info(r.URL.Path)

	u.UserMetrics.IncrementUsersCreated()
}

func (u UserHandler) ListUsers(rw http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(u.Options.ID()).Start(r.Context(), "users.ListUsers")
	defer span.End()

	u.Options.Logger().Info(r.URL.Path)

	time.Sleep(5 * time.Second)

	func() {
		_, span := otel.Tracer(u.Options.ID()).Start(ctx, "users.ListUsers.Wait")
		defer span.End()

		time.Sleep(5 * time.Second)
	}()

	u.UserMetrics.IncrementUsersCreated()
}

func (u UserHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	rid := uuid.NewString()
	ctx, span := otel.Tracer(u.Options.ID()).Start(r.Context(), "users.DeleteUser")
	defer span.End()

	u.Options.Logger().Infof("Request-ID: %v", rid)

	span.SetAttributes(attribute.Key("request_id").String(rid))

	u.Options.Logger().Info(r.URL.Path)

	time.Sleep(2 * time.Second)

	func() {
		_, span := otel.Tracer(u.Options.ID()).Start(ctx, "users.ListUsers.Wait")
		defer span.End()

		time.Sleep(3 * time.Second)
	}()

	time.Sleep(1 * time.Second)

	u.UserMetrics.IncrementUsersDeleted()
}

func (u UserHandler) GetUser(rw http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(u.Options.ID()).Start(r.Context(), "users.GetUser")
	defer span.End()

	time.Sleep(5 * time.Second)

	func() {
		_, span := otel.Tracer(u.Options.ID()).Start(ctx, "users.ListUsers.Wait")
		defer span.End()

		time.Sleep(5 * time.Second)
	}()
}

func (u UserHandler) UpdateUser(rw http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(u.Options.ID()).Start(r.Context(), "users.UpdateUser")
	defer span.End()

	time.Sleep(5 * time.Second)

	func() {
		_, span := otel.Tracer(u.Options.ID()).Start(ctx, "users.ListUsers.Wait")
		defer span.End()

		time.Sleep(5 * time.Second)
	}()
}

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
