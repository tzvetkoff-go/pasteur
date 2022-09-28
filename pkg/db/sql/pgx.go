//go:build !no_pgx
// +build !no_pgx

package sql

import (
	_ "github.com/jackc/pgx/v4"
)
