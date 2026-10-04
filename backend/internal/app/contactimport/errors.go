package contactimport

import (
	"errors"

	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

func asDomain(err error, target **domain.Error) bool { return errors.As(err, target) }
