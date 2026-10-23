package public

import (
	"github.com/gofiber/fiber/v2"

	httpx "github.com/Esca6585dev/habarchy/backend/internal/adapters/http"
	"github.com/Esca6585dev/habarchy/backend/internal/app/auth"
)

type phoneStartRequest struct {
	Phone string `json:"phone"`
}

type phoneVerifyRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

// phoneStart sends a login OTP to the phone (through the 0796 SMS provider).
func (h *Handlers) phoneStart(c *fiber.Ctx) error {
	var req phoneStartRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	res, err := h.PhoneAuth.Start(c.UserContext(), req.Phone, c.IP())
	if err != nil {
		return err
	}
	return httpx.OK(c, res)
}

// phoneVerify checks the code and returns JWT tokens, creating the user and
// their personal tenant on first login.
func (h *Handlers) phoneVerify(c *fiber.Ctx) error {
	var req phoneVerifyRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	sess := auth.Session{UserAgent: c.Get(fiber.HeaderUserAgent), IP: c.IP()}
	res, err := h.PhoneAuth.Verify(c.UserContext(), req.Phone, req.Code, sess)
	if err != nil {
		return err
	}
	return httpx.OK(c, fiber.Map{
		"tokens":  res.Tokens,
		"created": res.Created,
		"user":    fiber.Map{"id": res.User.ID, "phone": res.User.Phone, "full_name": res.User.FullName},
	})
}
