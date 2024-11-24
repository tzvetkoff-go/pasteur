package sql

import (
	dbSql "database/sql"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/tzvetkoff-go/logger"
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
func New(config *Config) (*SQL, error) {
	db, err := dbSql.Open(config.Driver, config.DSN)
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

	if config.AutoMigrate {
		err = result.Migrate()
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// Migrate ...
func (sql *SQL) Migrate() error {
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

// PrepareWithCache ...
func (sql *SQL) PrepareWithCache(q string) (*dbSql.Stmt, error) {
	if _, ok := sql.StatementCache[q]; !ok {
		stmt, err := sql.DB.Prepare(q)
		if err != nil {
			return nil, err
		}

		sql.Mux.Lock()
		sql.StatementCache[q] = stmt
		sql.Mux.Unlock()
	}

	return sql.StatementCache[q], nil
}

// Exec ...
func (sql *SQL) Exec(q string, args ...interface{}) (dbSql.Result, error) {
	stmt, err := sql.PrepareWithCache(q)
	if err != nil {
		return nil, err
	}

	logger.Debug("", logger.Fields{
		"query":      q,
		"query-args": args,
	})

	return stmt.Exec(args...)
}

// QueryRow ...
func (sql *SQL) QueryRow(q string, args ...interface{}) (*dbSql.Row, error) {
	stmt, err := sql.PrepareWithCache(q)
	if err != nil {
		return nil, err
	}

	logger.Debug("", logger.Fields{
		"query":      q,
		"query-args": args,
	})

	row := stmt.QueryRow(args...)
	return row, nil
}

// Query ...
func (sql *SQL) Query(q string, args ...interface{}) (*dbSql.Rows, error) {
	stmt, err := sql.PrepareWithCache(q)
	if err != nil {
		return nil, err
	}

	logger.Debug("", logger.Fields{
		"query":      q,
		"query-args": args,
	})

	return stmt.Query(args...)
}

// CreatePasteQuery ...
var CreatePasteQuery = stringutil.FormatQuery(`
	INSERT INTO pastes (private, filename, filetype, indent_style, indent_size, content, created_at)
				VALUES (?, ?, ?, ?, ?, ?, ?)
`)

// CreatePaste ...
func (sql *SQL) CreatePaste(paste *model.Paste) (*model.Paste, error) {
	paste.CreatedAt = time.Now()
	result, err := sql.Exec(
		CreatePasteQuery,
		paste.Private,
		paste.Filename,
		paste.Filetype,
		paste.IndentStyle,
		paste.IndentSize,
		paste.Content,
		paste.CreatedAt,
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

// GetPasteByIDQuery ...
var GetPasteByIDQuery = stringutil.FormatQuery(`
	SELECT id, private, indent_style, indent_size, filename, filetype, content, created_at
	  FROM pastes
	 WHERE id = ?
`)

// GetPasteByID ...
func (sql *SQL) GetPasteByID(id int) (*model.Paste, error) {
	row, err := sql.QueryRow(GetPasteByIDQuery, id)
	if err != nil {
		return nil, err
	}

	paste := &model.Paste{}
	err = row.Scan(
		&paste.ID,
		&paste.Private,
		&paste.IndentStyle,
		&paste.IndentSize,
		&paste.Filename,
		&paste.Filetype,
		&paste.Content,
		&paste.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return paste, nil
}

// CountPastesQuery ...
var CountPastesQuery = stringutil.FormatQuery(`
	SELECT COUNT(*)
	  FROM pastes
`)

// ListPastesQuery ...
var ListPastesQuery = stringutil.FormatQuery(`
	SELECT id, private, indent_style, indent_size, filename, filetype, content, created_at
	  FROM pastes
`)

// PaginatePastes ...
func (sql *SQL) PaginatePastes( // revive:disable-line:function-result-limit
	page int,
	perPage int,
	distance int,
	order string,
	conditions ...interface{},
) (
	*model.PaginatedPasteList,
	error,
) {
	countQuery := CountPastesQuery
	listQuery := ListPastesQuery
	stringConditions := []string{}
	sqlConditions := []interface{}{}
	offset := page*perPage - perPage

	if len(conditions) > 0 {
		for _, condition := range conditions {
			switch condition := condition.(type) {
			case map[string]interface{}:
				for key, value := range condition {
					stringConditions = append(stringConditions, key+" = ?")
					sqlConditions = append(sqlConditions, value)
				}
			case string:
				stringConditions = append(stringConditions, condition)
			default:
				sqlConditions = append(sqlConditions, condition)
			}
		}
	}

	if len(stringConditions) > 0 {
		countQuery += " WHERE " + strings.Join(stringConditions, " AND ")
		listQuery += " WHERE " + strings.Join(stringConditions, " AND ")
	}

	if order == "ASC" {
		listQuery += " ORDER BY id ASC"
	} else {
		listQuery += " ORDER BY id DESC"
	}

	listQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", perPage, offset)

	// SELECT COUNT(*) ...
	row, err := sql.QueryRow(countQuery, sqlConditions...)
	if err != nil {
		return nil, err
	}
	total := 0
	err = row.Scan(&total)
	if err != nil {
		return nil, err
	}
	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	// SELECT * ...
	rows, err := sql.Query(listQuery, sqlConditions...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := &model.PaginatedPasteList{
		Pastes: []*model.Paste{},
		Pagination: &model.Pagination{
			CurrentPage:         page,
			PaginationStartPage: int(math.Max(float64(page-distance), 2)),
			PaginationEndPage:   int(math.Min(float64(page+distance+1), float64(totalPages))),
			ItemsPerPage:        perPage,
			TotalItems:          total,
			TotalPages:          totalPages,
		},
	}
	for rows.Next() {
		paste := &model.Paste{}
		err = rows.Scan(
			&paste.ID,
			&paste.Private,
			&paste.IndentStyle,
			&paste.IndentSize,
			&paste.Filename,
			&paste.Filetype,
			&paste.Content,
			&paste.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		result.Pastes = append(result.Pastes, paste)
	}

	return result, nil
}
