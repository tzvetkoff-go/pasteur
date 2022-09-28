//go:build !no_mysql
// +build !no_mysql

package sql

import (
	_ "github.com/go-sql-driver/mysql"
)
