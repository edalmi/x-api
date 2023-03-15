package internal

import (
	"github.com/edalmi/x-api/internal/caching"
	"github.com/edalmi/x-api/internal/logging"
	"github.com/edalmi/x-api/internal/pubsub"
	"github.com/edalmi/x-api/internal/queue"
	"github.com/prometheus/client_golang/prometheus"
)

type Options struct {
	Cache   caching.Cache
	Logger  logging.Logger
	Pubsub  pubsub.Pubsub
	Queue   queue.Queue
	Metrics prometheus.Registerer
}
