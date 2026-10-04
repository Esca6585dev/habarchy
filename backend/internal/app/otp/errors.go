package otp

import "errors"

func errorsAs(err error, target **domainError) bool { return errors.As(err, target) }
