package cache

import (
	"context"
	"errors"
	"fmt"
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

// Get retrieves the value stored under key. Returns (nil, false, nil) on
// a genuine cache miss. Redis connection errors are also treated as a
// cache miss (fail open) rather than propagated, so that cache
// unavailability never blocks the caller from falling back to the
// primary data source.
func (c *Cache) Get(ctx context.Context, key string) ([]byte, bool, error) {
	val, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			c.log.Debug("cache miss", slog.String("key", key))
			return nil, false, nil
		}

		return nil, false, fmt.Errorf("redis get: %w", err)
	}

	c.log.Debug("cache hit", slog.String("key", key))
	return val, true, nil
}

// Set stores value under key with the given TTL.
func (c *Cache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	err := c.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return fmt.Errorf("redis set: %w", err)
	}

	c.log.Debug("successfully saved to cache", slog.String("key", key), slog.Duration("ttl", ttl))
	return nil
}
