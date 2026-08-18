// Package cache prepares Redis infrastructure: connection, health checks and
// graceful shutdown. No business-specific caching logic lives here.
package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/yourorg/go-fiber-template/internal/config"
)

// Cache wraps a go-redis client.
type Cache struct {
	client *redis.Client
}

// New connects to Redis.
func New(ctx context.Context, cfg config.Redis) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr(),
		Password:     cfg.Password,
		DB:           cfg.DB,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx).Err(); err != nil {
		return nil, fmt.Errorf("cache: ping: %w", err)
	}

	return &Cache{client: client}, nil
}

// Client exposes the underlying go-redis client.
func (c *Cache) Client() *redis.Client { return c.client }

// Ping verifies connectivity.
func (c *Cache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

// Close gracefully closes the client.
func (c *Cache) Close() error {
	if c.client == nil {
		return nil
	}
	return c.client.Close()
}
