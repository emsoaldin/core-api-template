package migrations

import "embed"

//go:embed *.sql
// Data holds all migration files embedded
var Data embed.FS
