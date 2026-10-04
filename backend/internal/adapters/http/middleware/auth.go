package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/Esca6585dev/habarchy/backend/internal/app/audit"
	"github.com/Esca6585dev/habarchy/backend/internal/app/auth"
	"github.com/Esca6585dev/habarchy/backend/internal/app/projects"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

// Locals keys.
const (
	LocalUserID = "user_id"
	LocalCaller = "api_caller"
)

// RequireJWT authenticates admin requests with "Authorization: Bearer".
// On success the user id is stored in locals and an audit.Actor in the
// user context.
func RequireJWT(svc *auth.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		h := c.Get(fiber.HeaderAuthorization)
		if !strings.HasPrefix(h, "Bearer ") {
			return domain.ErrUnauthorized.WithMessage("missing bearer token")
		}
		claims, err := svc.ParseAccess(strings.TrimSpace(h[len("Bearer "):]))
		if err != nil {
			return err
		}
		uid, err := claims.UserID()
		if err != nil {
			return domain.ErrUnauthorized
		}
		c.Locals(LocalUserID, uid)
		c.SetUserContext(audit.WithActor(c.UserContext(), audit.Actor{UserID: uid, IP: c.IP(), UserAgent: c.Get(fiber.HeaderUserAgent)}))
		return c.Next()
	}
}

// UserID returns the authenticated admin user id.
func UserID(c *fiber.Ctx) uuid.UUID {
	id, _ := c.Locals(LocalUserID).(uuid.UUID)
	return id
}

// APIKeyConfig tunes RequireAPIKey.
type APIKeyConfig struct {
	SignatureTolerance timeDuration
}

// RequireAPIKey authenticates public API requests with X-Api-Key and the
// optional X-Signature / X-Timestamp pair.
func RequireAPIKey(svc *projects.Service, cfg APIKeyConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		caller, err := svc.AuthenticateAPIKey(c.UserContext(), c.Get("X-Api-Key"), c.IP(), projects.SignatureInput{
			Signature: c.Get("X-Signature"),
			Timestamp: c.Get("X-Timestamp"),
			Method:    c.Method(),
			Path:      c.Path(),
			Body:      c.Body(),
			Tolerance: cfg.SignatureTolerance,
		})
		if err != nil {
			return err
		}
		c.Locals(LocalCaller, caller)
		return c.Next()
	}
}

// Caller returns the authenticated API client.
func Caller(c *fiber.Ctx) *projects.Caller {
	caller, _ := c.Locals(LocalCaller).(*projects.Caller)
	return caller
}

// RequireScope rejects keys lacking the scope.
func RequireScope(scope domain.APIKeyScope) fiber.Handler {
	return func(c *fiber.Ctx) error {
		caller := Caller(c)
		if caller == nil {
			return domain.ErrUnauthorized
		}
		if !caller.HasScope(scope) {
			return domain.ErrForbiddenScope.WithDetails(map[string]any{"required_scope": scope})
		}
		return c.Next()
	}
}
