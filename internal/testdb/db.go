// Package testdb prepares isolated databases for integration tests.
package testdb

import (
	"context"
	"database/sql"
	"io/fs"
	"path/filepath"
	"testing"

	"example.com/project-rest/internal/store"
	"example.com/project-rest/migrations"
	"github.com/pressly/goose/v3"
)

func New(t *testing.T, service string) *sql.DB {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), service+".db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	provider := Provider(t, db, service)
	if _, err := provider.Up(context.Background()); err != nil {
		t.Fatal(err)
	}
	return db
}

func Provider(t *testing.T, db *sql.DB, service string) *goose.Provider {
	t.Helper()
	files, err := fs.Sub(migrations.Files, service)
	if err != nil {
		t.Fatal(err)
	}
	provider, err := goose.NewProvider(goose.DialectSQLite3, db, files)
	if err != nil {
		t.Fatal(err)
	}
	return provider
}
