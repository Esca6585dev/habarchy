package admin

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

// contextWithCancel derives a context that ends when the stream writer
// returns. fasthttp's body stream writer runs after the handler has
// returned, so c.UserContext() must not be used inside it.
func contextWithCancel(_ *fiber.Ctx) (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}
