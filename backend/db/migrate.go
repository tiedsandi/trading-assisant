// Package db owns the versioned, embedded PostgreSQL migrations.
package db

import (
	"database/sql"
	"embed"
	"io/fs"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/lock"
)

//go:embed migrations/*.sql
var migrations embed.FS

// NewProvider shares the pgx driver with the API. Close the returned SQL handle
// before closing pool. A PostgreSQL advisory lock serializes migration runners.
func NewProvider(pool *pgxpool.Pool) (*goose.Provider, *sql.DB, error) {
	source, err := fs.Sub(migrations, "migrations")
	if err != nil {
		return nil, nil, err
	}
	locker, err := lock.NewPostgresSessionLocker()
	if err != nil {
		return nil, nil, err
	}
	sqlDB := stdlib.OpenDBFromPool(pool)
	provider, err := goose.NewProvider(goose.DialectPostgres, sqlDB, source,
		goose.WithSessionLocker(locker), goose.WithDisableGlobalRegistry(true))
	if err != nil {
		_ = sqlDB.Close()
		return nil, nil, err
	}
	return provider, sqlDB, nil
}
