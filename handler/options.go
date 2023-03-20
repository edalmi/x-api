package handler

import (
	"github.com/edalmi/x-api/caching"
	"github.com/edalmi/x-api/database"
	"github.com/edalmi/x-api/logging"
	"github.com/edalmi/x-api/pubsub"
	"github.com/edalmi/x-api/queue"
	"github.com/prometheus/client_golang/prometheus"
)

type HandlerOpts interface {
	Queue() queue.Queue
	Pubsub() pubsub.Pubsub
	Cache() caching.Cache
	Logger() logging.Logger
	Prometheus() prometheus.Registerer
	DB() *database.DB
	ID() string
}
