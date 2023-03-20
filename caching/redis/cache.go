package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewCache(rdb *redis.Client) *Cache {
	return &Cache{
		client: rdb,
	}
}

type Cache struct {
	client *redis.Client
}

func (c Cache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c Cache) Set(ctx context.Context, key string, value string, dur time.Duration) error {
	return c.client.Set(ctx, key, value, dur).Err()
}

func (c Cache) Close() error {
	return c.client.Close()
}
