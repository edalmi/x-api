package server

import (
	"github.com/edalmi/x-api/caching"
	"github.com/edalmi/x-api/database"
	"github.com/edalmi/x-api/logging"
	"github.com/edalmi/x-api/pubsub"
	"github.com/edalmi/x-api/queue"
	"github.com/prometheus/client_golang/prometheus"
)

type Options struct {
	db      *database.DB
	cache   caching.Cache
	logger  logging.Logger
	pubsub  pubsub.Pubsub
	queue   queue.Queue
	metrics prometheus.Registerer
}

func (o *Options) Cache() caching.Cache {
	return o.Cache
}

func (o *Options) Logger() logging.Logger {
	return o.logger
}

func (o *Options) Pubsub() pubsub.Pubsub {
	return o.pubsub
}

func (o *Options) Queue() queue.Queue {
	return o.queue
}

func (o *Options) Metrics() prometheus.Registerer {
	return o.metrics
}

func (o *Options) DB() *database.DB {
	return o.db
}
