package middleware

import (
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

// Limiter is the token-bucket interface the Redis client satisfies.
type Limiter interface {
	Allow(ctx context.Context, key string, rate float64, burst int) (bool, time.Duration, error)
}

// RateLimitByKey limits public API calls per API key. Burst is 2x rate.
func RateLimitByKey(l Limiter, perSec float64) fiber.Handler {
	burst := int(perSec * 2)
	if burst < 1 {
		burst = 1
	}
	return func(c *fiber.Ctx) error {
		k := Caller(c)
		if k == nil {
			return c.Next()
		}
		ok, wait, err := l.Allow(c.UserContext(), "api:"+k.Key.ID.String(), perSec, burst)
		if err != nil {
			return c.Next() // Redis trouble must not take the API down
		}
		if !ok {
			c.Set("Retry-After", strconv.Itoa(int(wait.Seconds())+1))
			return domain.ErrRateLimited.WithDetails(map[string]any{"retry_after_ms": wait.Milliseconds()})
		}
		return c.Next()
	}
}
