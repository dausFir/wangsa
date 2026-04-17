package database

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

//go:embed schema.sql
var embeddedSchema []byte

// NewPostgresDB opens a connection pool to PostgreSQL using DATABASE_URL.
// Supports both full DSN (postgres://user:pass@host/db) and individual params.
func NewPostgresDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// Connection pool tuning — sensible defaults for a web app
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return db, nil
}

// RunMigrations applies the schema (idempotent — uses IF NOT EXISTS throughout).
// If SCHEMA_PATH env is set that file is used; otherwise the embedded schema is used.
func RunMigrations(db *sqlx.DB, schemaPath string) error {
	var schema []byte
	if schemaPath != "" {
		b, err := os.ReadFile(schemaPath)
		if err != nil {
			return fmt.Errorf("read schema file '%s': %w", schemaPath, err)
		}
		schema = b
	} else {
		schema = embeddedSchema
	}
	if _, err := db.Exec(string(schema)); err != nil {
		return fmt.Errorf("execute migrations: %w", err)
	}
	return nil
}
