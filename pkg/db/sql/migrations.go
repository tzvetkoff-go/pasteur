package sql

import (
	"github.com/tzvetkoff-go/pasteur/pkg/stringutil"
)

// Migrations ...
var Migrations = []Migration{
	{
		Version: "0001",
		Migrate: func(sql *SQL) error {
			_, err := sql.Exec(stringutil.TrimHeredoc(`
				CREATE TABLE pastes (
					id           INTEGER PRIMARY KEY AUTOINCREMENT,
					indent_style VARCHAR(255),
					indent_size  VARCHAR(255),
					mime_type    VARCHAR(255),
					filename     VARCHAR(255),
					content      TEXT
				);
			`))

			return err
		},
	},
}
