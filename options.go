package xapi

import "github.com/prometheus/client_golang/prometheus"

type Options struct {
	Cache   Cache
	Logger  Logger
	Pubsub  Pubsub
	Queue   Queue
	Metrics prometheus.Registerer
}
