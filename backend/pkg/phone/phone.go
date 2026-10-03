// Package phone normalizes phone numbers to E.164.
package phone

import (
	"errors"
	"strings"

	"github.com/nyaruka/phonenumbers"
)

// DefaultRegion is applied to national-format numbers (no leading +).
// Turkmenistan is Habarchy's primary market.
const DefaultRegion = "TM"

// ErrInvalid is returned for numbers that cannot be parsed or are not
// valid for their region.
var ErrInvalid = errors.New("invalid phone number")

// Normalize parses raw and returns it in E.164 form (+99365123456).
// Spaces, dashes, dots and parentheses are tolerated. Numbers without a
// country code are interpreted in region (ISO 3166-1 alpha-2); an empty
// region falls back to DefaultRegion.
func Normalize(raw, region string) (string, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "", ErrInvalid
	}
	// Accept the common "00" international prefix.
	if strings.HasPrefix(s, "00") {
		s = "+" + s[2:]
	}
	if region == "" {
		region = DefaultRegion
	}
	num, err := phonenumbers.Parse(s, strings.ToUpper(region))
	if err != nil {
		return "", ErrInvalid
	}
	if !phonenumbers.IsValidNumber(num) {
		return "", ErrInvalid
	}
	return phonenumbers.Format(num, phonenumbers.E164), nil
}

// MustNormalize is Normalize for constants in tests and seeds; it panics on
// invalid input.
func MustNormalize(raw string) string {
	s, err := Normalize(raw, "")
	if err != nil {
		panic(err)
	}
	return s
}

// Region returns the ISO region of an E.164 number, or "" if unknown.
func Region(e164 string) string {
	num, err := phonenumbers.Parse(e164, "")
	if err != nil {
		return ""
	}
	return phonenumbers.GetRegionCodeForNumber(num)
}
