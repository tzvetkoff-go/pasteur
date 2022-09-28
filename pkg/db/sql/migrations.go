package sql

import (
	"github.com/tzvetkoff-go/pasteur/pkg/stringutil"
)

// Migrations ...
var Migrations = []Migration{
	{
		Version: "0001",
		Migrate: func(sql *SQL) error {
			_, err := sql.Exec(stringutil.FormatQuery(`
				CREATE TABLE pastes (
					id           INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
					private      TINYINT NOT NULL,
					filename     VARCHAR(255) NOT NULL,
					filetype     VARCHAR(255) NOT NULL,
					indent_style VARCHAR(255) NOT NULL,
					indent_size  INTEGER NOT NULL,
					content      TEXT NOT NULL,
					created_at   DATETIME NOT NULL
				);
			`))

			return err
		},
	},
}
