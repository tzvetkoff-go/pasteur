package db

import (
	"github.com/tzvetkoff-go/pasteur/pkg/db/sql"
)

// Config ...
type Config struct {
	Type string     `yaml:"type"`
	SQL  sql.Config `yaml:"sql"`
}
