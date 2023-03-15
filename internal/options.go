package internal

import (
	"github.com/edalmi/x-api/caching"
	"github.com/edalmi/x-api/logging"
	"github.com/edalmi/x-api/pubsub"
	"github.com/edalmi/x-api/queue"
	"github.com/prometheus/client_golang/prometheus"
)

type Options struct {
	Cache   caching.Cache
	Logger  logging.Logger
	Pubsub  pubsub.Pubsub
	Queue   queue.Queue
	Metrics prometheus.Registerer
}
