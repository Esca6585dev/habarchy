// Package postgres wires pgx, goose migrations and the sqlc-generated
// queries together.
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres/sqlcgen"
	"github.com/Esca6585dev/habarchy/backend/internal/config"
	"github.com/Esca6585dev/habarchy/backend/migrations"
)

// DB bundles the connection pool with the typed query layer.
type DB struct {
	Pool    *pgxpool.Pool
	Queries *sqlcgen.Queries
}

// Connect opens a pool using cfg and verifies connectivity.
func Connect(ctx context.Context, cfg config.Database) (*DB, error) {
	pcfg, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("postgres: parse url: %w", err)
	}
	pcfg.MaxConns = cfg.MaxConns
	pcfg.MinConns = cfg.MinConns
	pcfg.MaxConnLifetime = cfg.MaxConnLifetime
	pcfg.HealthCheckPeriod = 30 * time.Second
	pcfg.AfterConnect = registerTypes

	pool, err := pgxpool.NewWithConfig(ctx, pcfg)
	if err != nil {
		return nil, fmt.Errorf("postgres: connect: %w", err)
	}
	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("postgres: ping: %w", err)
	}
	return &DB{Pool: pool, Queries: sqlcgen.New(pool)}, nil
}

// customTypes are the Postgres enum types (and their array forms) that pgx
// must know about to scan them into Go slices. Plain enum columns work
// without registration; arrays of enums do not.
var customTypes = []string{"channel", "_channel"}

// registerTypes loads the custom types into the connection's type map. It
// tolerates missing types so a connection can be opened before migrations
// have run; Migrate resets the pool afterwards.
func registerTypes(ctx context.Context, conn *pgx.Conn) error {
	for _, name := range customTypes {
		t, err := conn.LoadType(ctx, name)
		if err != nil {
			return nil //nolint:nilerr // type does not exist yet; see Migrate
		}
		conn.TypeMap().RegisterType(t)
	}
	return nil
}

// Close releases the pool.
func (db *DB) Close() { db.Pool.Close() }

// Ping is used by /readyz.
func (db *DB) Ping(ctx context.Context) error { return db.Pool.Ping(ctx) }

// WithTx runs fn inside a transaction, committing on nil error.
func (db *DB) WithTx(ctx context.Context, fn func(q *sqlcgen.Queries) error) error {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if err := fn(db.Queries.WithTx(tx)); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// Migrate applies all pending goose migrations from the embedded FS.
func (db *DB) Migrate(ctx context.Context, log zerolog.Logger) error {
	err := runGoose(ctx, db.Pool, log, func(p *goose.Provider) error {
		results, err := p.Up(ctx)
		for _, r := range results {
			log.Info().Str("migration", r.Source.Path).Dur("took", r.Duration).Msg("migration applied")
		}
		return err
	})
	if err != nil {
		return err
	}
	// Connections opened before the schema existed have no custom types
	// registered; recycle them so AfterConnect runs again.
	db.Pool.Reset()
	return nil
}

// MigrateDownAll rolls back every migration. Only used by tests.
func (db *DB) MigrateDownAll(ctx context.Context, log zerolog.Logger) error {
	return runGoose(ctx, db.Pool, log, func(p *goose.Provider) error {
		_, err := p.DownTo(ctx, 0)
		return err
	})
}

// MigrationVersion returns the current schema version.
func (db *DB) MigrationVersion(ctx context.Context) (int64, error) {
	var v int64
	err := runGoose(ctx, db.Pool, zerolog.Nop(), func(p *goose.Provider) error {
		var err error
		v, err = p.GetDBVersion(ctx)
		return err
	})
	return v, err
}

func runGoose(_ context.Context, pool *pgxpool.Pool, _ zerolog.Logger, fn func(*goose.Provider) error) error {
	sqlDB := stdlib.OpenDBFromPool(pool)
	defer func(db *sql.DB) { _ = db.Close() }(sqlDB)
	provider, err := goose.NewProvider(goose.DialectPostgres, sqlDB, migrations.FS)
	if err != nil {
		return fmt.Errorf("goose: %w", err)
	}
	if err := fn(provider); err != nil {
		return fmt.Errorf("goose: %w", err)
	}
	return nil
}
