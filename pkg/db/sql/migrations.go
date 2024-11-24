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
			if err != nil {
				return err
			}

			_, err = sql.Exec(stringutil.FormatQuery(`
				CREATE INDEX index_pastes_on_filetype ON pastes (filetype);
			`))
			if err != nil {
				return err
			}

			_, err = sql.Exec(stringutil.FormatQuery(`
				INSERT INTO pastes (
					private, filename, filetype, indent_style, indent_size, content, created_at
				) VALUES (
					1,
					'hello-world.txt',
					'plain',
					'spaces',
					4,
					'Hello, world!',
					'1987-01-07 10:45:00.000'
				);
			`))

			return err
		},
	},
}
