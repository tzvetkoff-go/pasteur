package sql

import (
	dbSql "database/sql"
	"sync"

	"github.com/tzvetkoff-go/logger"
	"github.com/tzvetkoff-go/pasteur/pkg/config"
	"github.com/tzvetkoff-go/pasteur/pkg/model"
	"github.com/tzvetkoff-go/pasteur/pkg/stringutil"
)

// SQL ...
type SQL struct {
	DB             *dbSql.DB
	StatementCache map[string]*dbSql.Stmt
	Mux            sync.Mutex
}

// New ...
func New(sqlConfig *config.SQL) (*SQL, error) {
	db, err := dbSql.Open(sqlConfig.Driver, sqlConfig.DSN)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	result := &SQL{
		DB:             db,
		StatementCache: map[string]*dbSql.Stmt{},
	}

	err = result.Setup()
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Setup ...
func (sql *SQL) Setup() error {
	_, err := sql.DB.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (version VARCHAR(255))`)
	if err != nil {
		return err
	}

	rows, err := sql.DB.Query(`SELECT version FROM schema_migrations`)
	if err != nil {
		return err
	}
	defer rows.Close()

	schemaMigrations := []string{}
	for rows.Next() {
		version := ""
		err = rows.Scan(&version)
		if err != nil {
			return err
		}

		schemaMigrations = append(schemaMigrations, version)
	}

	for _, migration := range Migrations {
		migrated := false
		for _, version := range schemaMigrations {
			if version == migration.Version {
				migrated = true
				break
			}
		}

		if migrated {
			continue
		}

		logger.Info("Migrating...", logger.Fields{
			"version": migration.Version,
		})

		err = migration.Migrate(sql)
		if err != nil {
			logger.Error("Error", logger.Fields{
				"version": migration.Version,
				"error":   err,
			})

			return err
		}

		logger.Info("Migrated", logger.Fields{
			"version": migration.Version,
		})

		sql.Exec(`INSERT INTO schema_migrations (version) VALUES (?)`, migration.Version)
	}

	return nil
}

// Exec ...
func (sql *SQL) Exec(q string, args ...interface{}) (dbSql.Result, error) {
	if _, ok := sql.StatementCache[q]; !ok {
		stmt, err := sql.DB.Prepare(q)
		if err != nil {
			return nil, err
		}

		sql.Mux.Lock()
		sql.StatementCache[q] = stmt
		sql.Mux.Unlock()
	}

	logger.Debug("", logger.Fields{
		"query":      q,
		"query-args": args,
	})

	return sql.StatementCache[q].Exec(args...)
}

// QueryRow ...
func (sql *SQL) QueryRow(q string, args ...interface{}) (*dbSql.Row, error) {
	if _, ok := sql.StatementCache[q]; !ok {
		stmt, err := sql.DB.Prepare(q)
		if err != nil {
			return nil, err
		}

		sql.Mux.Lock()
		sql.StatementCache[q] = stmt
		sql.Mux.Unlock()
	}

	logger.Debug("", logger.Fields{
		"query":      q,
		"query-args": args,
	})

	row := sql.StatementCache[q].QueryRow(args...)
	return row, nil
}

// Query ...
func (sql *SQL) Query(q string, args ...interface{}) (*dbSql.Rows, error) {
	if _, ok := sql.StatementCache[q]; !ok {
		stmt, err := sql.DB.Prepare(q)
		if err != nil {
			return nil, err
		}

		sql.Mux.Lock()
		sql.StatementCache[q] = stmt
		sql.Mux.Unlock()
	}

	logger.Debug("", logger.Fields{
		"query":      q,
		"query-args": args,
	})

	return sql.StatementCache[q].Query(args...)
}

// CreatePaste ...
func (sql *SQL) CreatePaste(paste *model.Paste) (*model.Paste, error) {
	result, err := sql.Exec(
		stringutil.TrimHeredoc(`
			INSERT INTO pastes (indent_style, indent_size, mime_type, filename, content) VALUES (?, ?, ?, ?, ?)
		`),
		paste.IndentStyle,
		paste.IndentSize,
		paste.MimeType,
		paste.Filename,
		paste.Content,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	paste.ID = int(id)

	return paste, nil
}

// RetrievePasteByID ...
func (sql *SQL) RetrievePasteByID(id int) (*model.Paste, error) {
	row, err := sql.QueryRow(
		stringutil.TrimHeredoc(`
			SELECT id, indent_style, indent_size, mime_type, filename, content FROM pastes WHERE id = ?
		`),
	)
	if err != nil {
		return nil, err
	}

	result := &model.Paste{}
	err = row.Scan(
		&result.ID,
		&result.IndentStyle,
		&result.IndentSize,
		&result.MimeType,
		&result.Filename,
		&result.Content,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}
