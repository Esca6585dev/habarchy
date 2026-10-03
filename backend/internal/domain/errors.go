package domain

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrorCode is the machine-readable code returned in the API error envelope.
type ErrorCode string

const (
	CodeValidationFailed    ErrorCode = "validation_failed"
	CodeUnauthorized        ErrorCode = "unauthorized"
	CodeForbiddenScope      ErrorCode = "forbidden_scope"
	CodeNotFound            ErrorCode = "not_found"
	CodeQuotaExceeded       ErrorCode = "quota_exceeded"
	CodeRateLimited         ErrorCode = "rate_limited"
	CodeDuplicateRequest    ErrorCode = "duplicate_request"
	CodeProviderUnavailable ErrorCode = "provider_unavailable"
	CodeInvalidRecipient    ErrorCode = "invalid_recipient"
	CodeConflict            ErrorCode = "conflict"
	CodeInternal            ErrorCode = "internal_error"
)

// Sentinel errors. Match with errors.Is; wrap with fmt.Errorf("%w: ...")
// to add context without losing the code.
var (
	ErrValidation          = &Error{Code: CodeValidationFailed, Message: "validation failed"}
	ErrUnauthorized        = &Error{Code: CodeUnauthorized, Message: "unauthorized"}
	ErrForbiddenScope      = &Error{Code: CodeForbiddenScope, Message: "api key lacks the required scope"}
	ErrNotFound            = &Error{Code: CodeNotFound, Message: "not found"}
	ErrQuotaExceeded       = &Error{Code: CodeQuotaExceeded, Message: "project quota exceeded"}
	ErrRateLimited         = &Error{Code: CodeRateLimited, Message: "rate limited"}
	ErrDuplicateRequest    = &Error{Code: CodeDuplicateRequest, Message: "duplicate request"}
	ErrProviderUnavailable = &Error{Code: CodeProviderUnavailable, Message: "no provider available"}
	ErrInvalidRecipient    = &Error{Code: CodeInvalidRecipient, Message: "invalid recipient"}
	ErrConflict            = &Error{Code: CodeConflict, Message: "conflict"}
	ErrInternal            = &Error{Code: CodeInternal, Message: "internal error"}
)

// Error is a domain error carrying an API code and optional details.
type Error struct {
	Code    ErrorCode
	Message string
	Details map[string]any
}

func (e *Error) Error() string { return fmt.Sprintf("%s: %s", e.Code, e.Message) }

// Is makes sentinel matching work on wrapped and derived errors: two domain
// errors are equal when their codes match.
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	return ok && t.Code == e.Code
}

// WithDetails returns a copy of e with details attached.
func (e *Error) WithDetails(details map[string]any) *Error {
	return &Error{Code: e.Code, Message: e.Message, Details: details}
}

// WithMessage returns a copy of e with a more specific message.
func (e *Error) WithMessage(msg string) *Error {
	return &Error{Code: e.Code, Message: msg, Details: e.Details}
}

// CodeOf extracts the ErrorCode from any error, defaulting to internal_error.
func CodeOf(err error) ErrorCode {
	var de *Error
	if errors.As(err, &de) {
		return de.Code
	}
	return CodeInternal
}

// HTTPStatus maps an error code to an HTTP status.
func (c ErrorCode) HTTPStatus() int {
	switch c {
	case CodeValidationFailed, CodeInvalidRecipient:
		return http.StatusBadRequest
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbiddenScope:
		return http.StatusForbidden
	case CodeNotFound:
		return http.StatusNotFound
	case CodeConflict, CodeDuplicateRequest:
		return http.StatusConflict
	case CodeQuotaExceeded, CodeRateLimited:
		return http.StatusTooManyRequests
	case CodeProviderUnavailable:
		return http.StatusServiceUnavailable
	}
	return http.StatusInternalServerError
}
