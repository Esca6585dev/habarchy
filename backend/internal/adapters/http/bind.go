package http

import (
	"errors"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

// Validate is the shared validator instance. Struct tags use the
// go-playground syntax: `validate:"required,email"`.
var Validate = validator.New(validator.WithRequiredStructEnabled())

func init() {
	// Use json tag names in error details so clients see "daily_quota",
	// not "DailyQuota".
	Validate.RegisterTagNameFunc(func(f reflectStructField) string {
		name := strings.SplitN(f.Tag.Get("json"), ",", 2)[0]
		if name == "-" || name == "" {
			return f.Name
		}
		return name
	})
}

// Bind parses the JSON body into dst and validates it. Errors are domain
// validation errors with a details map of field -> rule.
func Bind(c *fiber.Ctx, dst any) error {
	if len(c.Body()) == 0 {
		return domain.ErrValidation.WithMessage("request body is required")
	}
	if err := c.BodyParser(dst); err != nil {
		return domain.ErrValidation.WithMessage("malformed JSON body")
	}
	return ValidateStruct(dst)
}

// ValidateStruct runs the validator and converts its errors.
func ValidateStruct(v any) error {
	err := Validate.Struct(v)
	if err == nil {
		return nil
	}
	var verrs validator.ValidationErrors
	if !errors.As(err, &verrs) {
		return domain.ErrValidation
	}
	details := make(map[string]any, len(verrs))
	for _, fe := range verrs {
		rule := fe.Tag()
		if fe.Param() != "" {
			rule += "=" + fe.Param()
		}
		details[fe.Field()] = rule
	}
	return domain.ErrValidation.WithDetails(details)
}

// ParamUUID reads a UUID path parameter.
func ParamUUID(c *fiber.Ctx, name string) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Params(name))
	if err != nil {
		return uuid.Nil, domain.ErrValidation.WithDetails(map[string]any{name: "must be a uuid"})
	}
	return id, nil
}

// Page holds offset pagination parsed from ?limit=&offset=.
type Page struct {
	Limit  int32
	Offset int32
}

// ParsePage reads pagination with bounds [1, max].
func ParsePage(c *fiber.Ctx, def, maxLimit int) Page {
	limit, _ := strconv.Atoi(c.Query("limit"))
	if limit <= 0 {
		limit = def
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	offset, _ := strconv.Atoi(c.Query("offset"))
	if offset < 0 {
		offset = 0
	}
	return Page{Limit: int32(limit), Offset: int32(offset)} //nolint:gosec // bounded above
}

// Created writes a 201 envelope.
func Created(c *fiber.Ctx, data any) error { return JSON(c, fiber.StatusCreated, data, nil) }

// NoContent writes 204.
func NoContent(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) }
