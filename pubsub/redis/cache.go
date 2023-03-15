package redis

import "github.com/redis/go-redis/v9"

func NewCache(rdb *redis.Client) *Cache {
	return &Cache{
		client: rdb,
	}
}

type Cache struct {
	client *redis.Client
}
