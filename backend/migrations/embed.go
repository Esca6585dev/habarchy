// Package migrations embeds the goose SQL migrations so binaries can apply
// them without the files being present on disk.
package migrations

import "embed"

// FS contains every *.sql migration.
//
//go:embed *.sql
var FS embed.FS
