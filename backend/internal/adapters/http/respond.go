// Package http contains the Fiber server, middleware and JSON envelope
// helpers shared by the public and admin APIs.
package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

// Envelope is the JSON shape of every response:
//
//	{ "data": ..., "meta": {...}, "error": {"code","message","details"} }
type Envelope struct {
	Data  any       `json:"data,omitempty"`
	Meta  any       `json:"meta,omitempty"`
	Error *APIError `json:"error,omitempty"`
}

// APIError is the error part of the envelope.
type APIError struct {
	Code    domain.ErrorCode `json:"code"`
	Message string           `json:"message"`
	Details map[string]any   `json:"details,omitempty"`
}

// OK writes a 200 envelope.
func OK(c *fiber.Ctx, data any) error { return JSON(c, fiber.StatusOK, data, nil) }

// JSON writes an envelope with an explicit status and optional meta.
func JSON(c *fiber.Ctx, status int, data, meta any) error {
	return c.Status(status).JSON(Envelope{Data: data, Meta: meta})
}

// Fail writes an error envelope derived from err. Domain errors keep their
// code and HTTP status; anything else becomes a 500 internal_error whose
// details are logged, not leaked.
func Fail(c *fiber.Ctx, err error) error {
	var de *domain.Error
	if errors.As(err, &de) {
		return c.Status(de.Code.HTTPStatus()).JSON(Envelope{Error: &APIError{
			Code: de.Code, Message: de.Message, Details: de.Details,
		}})
	}
	var fe *fiber.Error
	if errors.As(err, &fe) {
		code := domain.CodeInternal
		switch fe.Code {
		case fiber.StatusNotFound:
			code = domain.CodeNotFound
		case fiber.StatusBadRequest, fiber.StatusUnprocessableEntity, fiber.StatusRequestEntityTooLarge:
			code = domain.CodeValidationFailed
		case fiber.StatusUnauthorized:
			code = domain.CodeUnauthorized
		case fiber.StatusMethodNotAllowed:
			return c.Status(fe.Code).JSON(Envelope{Error: &APIError{Code: domain.CodeValidationFailed, Message: fe.Message}})
		}
		return c.Status(fe.Code).JSON(Envelope{Error: &APIError{Code: code, Message: fe.Message}})
	}
	log.Ctx(c.UserContext()).Error().Err(err).Str("path", c.Path()).Msg("unhandled error")
	return c.Status(fiber.StatusInternalServerError).JSON(Envelope{Error: &APIError{
		Code: domain.CodeInternal, Message: "internal error",
	}})
}
