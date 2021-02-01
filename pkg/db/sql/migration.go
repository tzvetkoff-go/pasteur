package sql

// Migration ...
type Migration struct {
	Version  string
	Migrate  func(*SQL) error
	Rollback func(*SQL) error
}
