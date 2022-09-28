package sql

// Config ...
type Config struct {
	Driver      string `yaml:"driver"`
	DSN         string `yaml:"dsn"`
	AutoMigrate bool   `yaml:"auto-migrate"`
}
