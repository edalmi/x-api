package memcached

import (
	"context"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

func New(addr []string) (*Cache, error) {
	return &Cache{
		client: memcache.New(addr...),
	}, nil
}

type Cache struct {
	client *memcache.Client
}

func (c Cache) Set(ctx context.Context, key, value string, expiration time.Duration) error {
	return c.client.Set(&memcache.Item{
		Key:        key,
		Value:      []byte(value),
		Expiration: int32(expiration),
	})
}

func (c Cache) Get(ctx context.Context, key string) (string, error) {
	value, err := c.client.Get(key)
	if err != nil {
		return "", err
	}

	return string(value.Value), nil
}

func (c Cache) Close() error {
	return c.client.Close()
}
