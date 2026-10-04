package admin

import (
	"github.com/gofiber/fiber/v2"

	httpx "github.com/Esca6585dev/habarchy/backend/internal/adapters/http"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/http/middleware"
	"github.com/Esca6585dev/habarchy/backend/internal/app/auth"
)

func session(c *fiber.Ctx) auth.Session {
	return auth.Session{UserAgent: c.Get(fiber.HeaderUserAgent), IP: c.IP()}
}

func (h *Handlers) login(c *fiber.Ctx) error {
	var req loginRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	tokens, user, err := h.Auth.Login(c.UserContext(), req.Email, req.Password, req.TOTPCode, session(c))
	if err != nil {
		return err
	}
	return httpx.OK(c, fiber.Map{"tokens": tokens, "user": toUser(user)})
}

func (h *Handlers) refresh(c *fiber.Ctx) error {
	var req refreshRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	tokens, err := h.Auth.Refresh(c.UserContext(), req.RefreshToken, session(c))
	if err != nil {
		return err
	}
	return httpx.OK(c, fiber.Map{"tokens": tokens})
}

func (h *Handlers) logout(c *fiber.Ctx) error {
	var req refreshRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	if err := h.Auth.Logout(c.UserContext(), req.RefreshToken); err != nil {
		return err
	}
	return httpx.NoContent(c)
}

func (h *Handlers) logoutAll(c *fiber.Ctx) error {
	if err := h.Auth.LogoutAll(c.UserContext(), middleware.UserID(c)); err != nil {
		return err
	}
	return httpx.NoContent(c)
}

func (h *Handlers) me(c *fiber.Ctx) error {
	user, err := h.Auth.User(c.UserContext(), middleware.UserID(c))
	if err != nil {
		return err
	}
	return httpx.OK(c, toUser(user))
}

func (h *Handlers) changePassword(c *fiber.Ctx) error {
	var req changePasswordRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	if err := h.Auth.ChangePassword(c.UserContext(), middleware.UserID(c), req.CurrentPassword, req.NewPassword); err != nil {
		return err
	}
	return httpx.NoContent(c)
}

func (h *Handlers) totpSetup(c *fiber.Ctx) error {
	setup, err := h.Auth.SetupTOTP(c.UserContext(), middleware.UserID(c))
	if err != nil {
		return err
	}
	return httpx.OK(c, setup)
}

func (h *Handlers) totpConfirm(c *fiber.Ctx) error {
	var req totpCodeRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	if err := h.Auth.ConfirmTOTP(c.UserContext(), middleware.UserID(c), req.Code); err != nil {
		return err
	}
	return httpx.NoContent(c)
}

func (h *Handlers) totpDisable(c *fiber.Ctx) error {
	var req passwordRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	if err := h.Auth.DisableTOTP(c.UserContext(), middleware.UserID(c), req.Password); err != nil {
		return err
	}
	return httpx.NoContent(c)
}
