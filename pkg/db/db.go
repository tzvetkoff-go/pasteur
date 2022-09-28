package db

import (
	"github.com/tzvetkoff-go/errors"

	"github.com/tzvetkoff-go/pasteur/pkg/db/sql"
	"github.com/tzvetkoff-go/pasteur/pkg/model"
)

// DefaultPerPage ...
const DefaultPerPage = 20

// DB ...
type DB interface {
	Migrate() error
	GetPasteByID(int) (*model.Paste, error)
	CreatePaste(*model.Paste) (*model.Paste, error)
	PaginatePastes(
		page int,
		perPage int,
		distance int,
		order string,
		conditions ...interface{},
	) (*model.PaginatedPasteList, error)
}

// New ...
func New(config *Config) (DB, error) {
	switch config.Type {
	case "sql":
		return sql.New(&config.SQL)
	default:
		return nil, errors.New("unknown database storage: %s", config.Type)
	}
}
