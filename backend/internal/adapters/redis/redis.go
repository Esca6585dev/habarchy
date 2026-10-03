// Package redis creates the shared Redis client used for queues, rate
// limits, idempotency and the OTP store.
package redis

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"github.com/Esca6585dev/habarchy/backend/internal/config"
)

// Client wraps go-redis so callers do not import it directly.
type Client struct {
	c    *goredis.Client
	opts *goredis.Options
}

// Connect parses cfg.URL (redis://user:pass@host:port/db) and pings.
func Connect(ctx context.Context, cfg config.Redis) (*Client, error) {
	opts, err := goredis.ParseURL(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("redis: parse url: %w", err)
	}
	c := goredis.NewClient(opts)
	client := &Client{c: c, opts: opts}
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx); err != nil {
		_ = c.Close()
		return nil, fmt.Errorf("redis: ping: %w", err)
	}
	return client, nil
}

// Ping checks connectivity; used by /readyz.
func (c *Client) Ping(ctx context.Context) error { return c.c.Ping(ctx).Err() }

// Close closes the underlying connection pool.
func (c *Client) Close() error { return c.c.Close() }

// Raw returns the underlying go-redis client for adapters that need it.
func (c *Client) Raw() *goredis.Client { return c.c }

// Options exposes the parsed connection options for libraries (asynq) that
// open their own connections.
func (c *Client) Options() *goredis.Options { return c.opts }
