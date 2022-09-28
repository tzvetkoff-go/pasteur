//go:build !no_sqlite3
// +build !no_sqlite3

package sql

import (
	_ "github.com/mattn/go-sqlite3"
)
