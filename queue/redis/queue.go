package redis

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/redis/go-redis/v9"

	"github.com/edalmi/x-api/queue"
)

type Redis struct {
	client *redis.Client
}

func (r *Redis) Push(ctx context.Context, name string, msg queue.Message) error {
	type data struct {
		ContentType string `json:"contentType"`
		Body        string `json:"body"`
	}

	d := data{
		ContentType: msg.ContentType,
		Body:        string(msg.Body),
	}

	dd, err := json.Marshal(d)
	if err != nil {
		return err
	}

	if err := r.client.RPush(ctx, name, string(dd)).Err(); err != nil {
		return err
	}

	return nil
}

func (r *Redis) Pull(ctx context.Context, name string) (<-chan *queue.Message, error) {
	return nil, errors.New("error")

}

func (r *Redis) Close() error {
	return errors.New("error")
}
