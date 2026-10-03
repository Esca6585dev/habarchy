// Package ids generates identifiers.
package ids

import "github.com/google/uuid"

// New returns a time-ordered UUID v7. Messages use v7 so the primary key
// index stays append-only and ids sort by creation time.
func New() uuid.UUID {
	id, err := uuid.NewV7()
	if err != nil {
		// NewV7 only fails when the OS random source is broken; fall back
		// to v4 rather than crash a request.
		return uuid.New()
	}
	return id
}

// Parse validates a UUID string.
func Parse(s string) (uuid.UUID, error) { return uuid.Parse(s) }
