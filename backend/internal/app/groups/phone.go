package groups

import "github.com/Esca6585dev/habarchy/backend/pkg/phone"

func normalizePhone(raw string) string {
	p, err := phone.Normalize(raw, "")
	if err != nil {
		return raw
	}
	return p
}
