package cache

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	log    *slog.Logger
	client *redis.Client
}

func NewCache(client *redis.Client, log *slog.Logger) *Cache {
	return &Cache{log: log, client: client}
}

func (c *Cache) Get(ctx context.Context, key string) ([]byte, bool, error) {
	val, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			c.log.Debug("cache miss", slog.String("key", key))
			return nil, false, nil
		}

		c.log.Error("redis connection error during GET", slog.String("key", key), slog.Any("error", err))
		return nil, false, nil
	}

	c.log.Debug("cache hit", slog.String("key", key))
	return val, true, nil
}

func (c *Cache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	err := c.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		c.log.Error("redis connection error during SET", slog.String("key", key), slog.Any("error", err))
		return nil
	}

	c.log.Debug("successfully saved to cache", slog.String("key", key), slog.Duration("ttl", ttl))
	return nil
}
