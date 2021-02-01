package db

import (
	"github.com/tzvetkoff-go/errors"

	"github.com/tzvetkoff-go/pasteur/pkg/config"
	"github.com/tzvetkoff-go/pasteur/pkg/db/filesystem"
	"github.com/tzvetkoff-go/pasteur/pkg/db/sql"
	"github.com/tzvetkoff-go/pasteur/pkg/model"
)

// DB ...
type DB interface {
	CreatePaste(*model.Paste) (*model.Paste, error)
	RetrievePasteByID(int) (*model.Paste, error)
}

// New ...
func New(dbConfig *config.DB) (DB, error) {
	switch dbConfig.Storage {
	case "filesystem":
		return filesystem.New(&dbConfig.FileSystem)
	case "sql":
		return sql.New(&dbConfig.SQL)
	default:
		return nil, errors.New("unknown database storage: %s", dbConfig.Storage)
	}
}
