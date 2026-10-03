// Package logger configures the process-wide zerolog logger.
package logger

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// New returns a structured JSON logger (or a human-readable console logger
// when pretty is true) at the given level, tagged with the service name.
// It also installs the logger as zerolog's global so library code using
// log.Ctx / log.Logger shares the configuration.
func New(service, level string, pretty bool) zerolog.Logger {
	lvl, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil || lvl == zerolog.NoLevel {
		lvl = zerolog.InfoLevel
	}
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.DurationFieldUnit = time.Millisecond

	var w io.Writer = os.Stdout
	if pretty {
		w = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.Kitchen}
	}
	l := zerolog.New(w).Level(lvl).With().Timestamp().Str("service", service).Logger()
	log.Logger = l
	return l
}
