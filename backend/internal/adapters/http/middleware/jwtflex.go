package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/Esca6585dev/habarchy/backend/internal/app/audit"
	"github.com/Esca6585dev/habarchy/backend/internal/app/auth"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

// AccessCookie is the cookie name the web admin uses for the access token.
const AccessCookie = "habarchy_access"

// RequireJWTFlexible accepts the access token from the Authorization
// header, the ?access_token= query parameter or the access cookie. Used for
// SSE (EventSource cannot set headers) and the asynqmon UI.
func RequireJWTFlexible(svc *auth.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := ""
		if h := c.Get(fiber.HeaderAuthorization); strings.HasPrefix(h, "Bearer ") {
			token = strings.TrimSpace(h[len("Bearer "):])
		}
		if token == "" {
			token = c.Query("access_token")
		}
		if token == "" {
			token = c.Cookies(AccessCookie)
		}
		if token == "" {
			return domain.ErrUnauthorized.WithMessage("missing access token")
		}
		claims, err := svc.ParseAccess(token)
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
