package http

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/http/middleware"
	"github.com/Esca6585dev/habarchy/backend/internal/config"
)

// Pinger is anything whose liveness /readyz should check.
type Pinger interface {
	Ping(ctx context.Context) error
}

// Deps are the external dependencies the server reports on and hands to
// handlers. Fields are added as steps of the build add features.
type Deps struct {
	DB    Pinger
	Redis Pinger
}

// NewServer builds the Fiber application with base middleware and health
// endpoints. Feature routers (public API, admin API) register on it.
func NewServer(cfg *config.Config, log zerolog.Logger, deps Deps) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:                 "habarchy",
		ReadTimeout:             cfg.HTTP.ReadTimeout,
		WriteTimeout:            cfg.HTTP.WriteTimeout,
		IdleTimeout:             cfg.HTTP.IdleTimeout,
		BodyLimit:               cfg.HTTP.BodyLimitBytes,
		DisableStartupMessage:   true,
		ProxyHeader:             proxyHeader(cfg.HTTP.TrustedProxies),
		EnableTrustedProxyCheck: len(cfg.HTTP.TrustedProxies) > 0,
		TrustedProxies:          cfg.HTTP.TrustedProxies,
		ErrorHandler:            Fail,
	})

	app.Use(recover.New(recover.Config{EnableStackTrace: !cfg.IsProduction()}))
	app.Use(middleware.RequestID())
	app.Use(middleware.Logger(log))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(cfg.HTTP.CORSOrigins, ","),
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Api-Key, X-Signature, X-Timestamp, X-Request-Id",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowCredentials: true,
		ExposeHeaders:    middleware.RequestIDHeader,
	}))

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return OK(c, fiber.Map{"status": "ok"})
	})
	app.Get("/readyz", readyz(deps))

	return app
}

func proxyHeader(trusted []string) string {
	if len(trusted) == 0 {
		return ""
	}
	return fiber.HeaderXForwardedFor
}

func readyz(deps Deps) fiber.Handler {
	return func(c *fiber.Ctx) error {
		checks := map[string]string{}
		healthy := true
		check := func(name string, p Pinger) {
			if p == nil {
				checks[name] = "skipped"
				return
			}
			if err := p.Ping(c.UserContext()); err != nil {
				checks[name] = "down"
				healthy = false
				return
			}
			checks[name] = "ok"
		}
		check("postgres", deps.DB)
		check("redis", deps.Redis)
		status := fiber.StatusOK
		if !healthy {
			status = fiber.StatusServiceUnavailable
		}
		return JSON(c, status, fiber.Map{"status": map[bool]string{true: "ready", false: "degraded"}[healthy], "checks": checks}, nil)
	}
}
