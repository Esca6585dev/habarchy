package postgres

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"

	"github.com/Esca6585dev/habarchy/backend/internal/config"
)

// testDB connects to HABARCHY_TEST_DATABASE_URL, applies all migrations and
// registers a full rollback as cleanup, so each test file gets a clean
// schema. Tests are skipped when the variable is unset.
func testDB(t *testing.T) *DB {
	t.Helper()
	url := os.Getenv("HABARCHY_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("HABARCHY_TEST_DATABASE_URL not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	t.Cleanup(cancel)

	db, err := Connect(ctx, config.Database{URL: url, MaxConns: 4, MinConns: 0, MaxConnLifetime: time.Minute})
	if err != nil {
		t.Fatal(err)
	}
	log := zerolog.Nop()
	// Start from zero so a previous aborted run cannot leak state.
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
