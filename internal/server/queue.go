package server

import (
	"errors"

	"github.com/edalmi/x-api/internal/config"
	"github.com/edalmi/x-api/internal/queue"
)

func setupQueue(cfg *config.Queue) (queue.Queue, error) {
	return nil, errors.New("error")
}
