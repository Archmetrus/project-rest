package migrations_test

import (
	"context"
	"path/filepath"
	"testing"

	"example.com/project-rest/internal/store"
	"example.com/project-rest/internal/testdb"
)

func TestAddressMigration(t *testing.T) {
	ctx := context.Background()
	db, err := store.Open(filepath.Join(t.TempDir(), "users.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	p := testdb.Provider(t, db, "user")
	if _, err := p.UpTo(ctx, 1); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("SELECT address FROM users"); err == nil {
		t.Fatal("initial schema must not contain address")
	}
	if _, err := db.Exec("INSERT INTO users (name, email) VALUES ('Ada', 'ada@example.com')"); err != nil {
		t.Fatal(err)
	}
	if _, err := p.Up(ctx); err != nil {
		t.Fatal(err)
	}
	var address string
	if err := db.QueryRow("SELECT address FROM users WHERE id = 1").Scan(&address); err != nil || address != "" {
		t.Fatalf("address default: %q, %v", address, err)
	}
	if _, err := p.Status(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := p.Down(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("SELECT address FROM users"); err == nil {
		t.Fatal("down must remove address")
	}
	if _, err := p.Up(ctx); err != nil {
		t.Fatal(err)
	}
	var name string
	if err := db.QueryRow("SELECT name FROM users WHERE id = 1").Scan(&name); err != nil || name != "Ada" {
		t.Fatalf("existing row lost: %q, %v", name, err)
	}
}

func TestHRMigration(t *testing.T) {
	db := testdb.New(t, "hr")
	p := testdb.Provider(t, db, "hr")
	ctx := context.Background()
	if _, err := p.Status(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := p.Down(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("SELECT * FROM leaves"); err == nil {
		t.Fatal("down must drop leaves")
	}
	if _, err := p.Up(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("INSERT INTO leaves (user_id, start_date, end_date) VALUES (999, '2026-10-20', '2026-10-01')"); err != nil {
		t.Fatal(err)
	}
}
