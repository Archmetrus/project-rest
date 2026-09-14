package migrations

import "embed"

// Files contains the service-owned schema migrations.
//
//go:embed user/*.sql hr/*.sql
var Files embed.FS
