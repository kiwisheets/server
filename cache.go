package server

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type Cache struct {
	client    redis.UniversalClient
	ttl       time.Duration
	appPrefix string
}

const apqPrefix = "apq:"

func newCache(redisAddress string, appPrefix string, ttl time.Duration) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddress,
	})

	err := client.Ping().Err()
	if err != nil {
		return nil, fmt.Errorf("could not create cache: %w", err)
	}

	return &Cache{
		client:    client,
		ttl:       ttl,
		appPrefix: appPrefix + ":",
	}, nil
}

func (c *Cache) Add(ctx context.Context, key string, value interface{}) {
	c.client.Set(buildKey(c.appPrefix, key), value, c.ttl)
}

func (c *Cache) Get(ctx context.Context, key string) (interface{}, bool) {
	s, err := c.client.Get(buildKey(c.appPrefix, key)).Result()
	if err != nil {
		return struct{}{}, false
	}
	return s, true
}

func buildKey(appPrefix, key string) string {
	return apqPrefix + appPrefix + key
}
