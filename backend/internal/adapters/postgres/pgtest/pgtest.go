// Package pgtest opens a migrated test database for integration tests.
package pgtest

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres"
	"github.com/Esca6585dev/habarchy/backend/internal/config"
)

// EnvVar names the connection string integration tests use.
const EnvVar = "HABARCHY_TEST_DATABASE_URL"

// lockID is the advisory lock key that serialises test packages sharing
// one database (go test runs packages in parallel processes).
const lockID = 0x4841_4241 // "HABA"

// Open connects to EnvVar, resets the schema (down then up) and registers
// a cleanup that rolls everything back. The test is skipped when EnvVar is
// unset. A session-level advisory lock is held for the duration of the
// test so packages running in parallel take turns.
func Open(t *testing.T) *postgres.DB {
	t.Helper()
	url := os.Getenv(EnvVar)
	if url == "" {
		t.Skip(EnvVar + " not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	t.Cleanup(cancel)

	lockConn, err := pgx.Connect(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := lockConn.Exec(ctx, "SELECT pg_advisory_lock($1)", lockID); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = lockConn.Exec(context.Background(), "SELECT pg_advisory_unlock($1)", lockID)
		_ = lockConn.Close(context.Background())
	})

	db, err := postgres.Connect(ctx, config.Database{URL: url, MaxConns: 4, MaxConnLifetime: time.Minute})
	if err != nil {
		t.Fatal(err)
	}
	log := zerolog.Nop()
	if err := db.MigrateDownAll(ctx, log); err != nil {
		t.Fatal(err)
	}
	if err := db.Migrate(ctx, log); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := db.MigrateDownAll(context.Background(), log); err != nil {
			t.Errorf("down migrations: %v", err)
		}
		db.Close()
	})
	return db
}
