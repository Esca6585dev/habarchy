// Package middleware holds Fiber middleware shared by all routers.
package middleware

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/rs/zerolog"

	"github.com/Esca6585dev/habarchy/backend/internal/domain"
	"github.com/Esca6585dev/habarchy/backend/pkg/ids"
)

// RequestIDHeader is echoed on every response.
const RequestIDHeader = "X-Request-Id"

// RequestID assigns a UUID v7 request id unless the client supplied one.
func RequestID() fiber.Handler {
	return requestid.New(requestid.Config{
		Header:    RequestIDHeader,
		Generator: func() string { return ids.New().String() },
	})
}

// Logger attaches a request-scoped zerolog logger (with request_id) to the
// user context and writes one structured line per request.
func Logger(base zerolog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		rid, _ := c.Locals(requestid.ConfigDefault.ContextKey).(string)
		l := base.With().Str("request_id", rid).Logger()
		c.SetUserContext(l.WithContext(c.UserContext()))

		err := c.Next()

		status := statusOf(c, err)
		ev := l.Info()
		switch {
		case status >= 500:
			ev = l.Error()
		case status >= 400:
			ev = l.Warn()
		}
		if err != nil {
			ev = ev.Err(err)
		}
		ev.Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status", status).
			Dur("duration", time.Since(start)).
			Str("ip", c.IP()).
			Msg("request")
		return err
	}
}

// statusOf returns the status the client will receive. When a handler
// returns an error the app's ErrorHandler has not run yet, so the status
// must be derived from the error the same way Fail does.
func statusOf(c *fiber.Ctx, err error) int {
	if err == nil {
		return c.Response().StatusCode()
	}
	var de *domain.Error
	if errors.As(err, &de) {
		return de.Code.HTTPStatus()
	}
	var fe *fiber.Error
	if errors.As(err, &fe) {
		return fe.Code
	}
	return fiber.StatusInternalServerError
}
