package logger

// Config ...
type Config struct {
	Backends []Backend `yaml:"backends"`
	Level    string    `yaml:"level"`
}

// Backend ...
type Backend struct {
	Type   string        `yaml:"type"`
	File   BackendFile   `yaml:"file"`
	Syslog BackendSyslog `yaml:"syslog"`
}

// BackendFile ...
type BackendFile struct {
	Destination string `yaml:"destination"`
}

// BackendSyslog ...
type BackendSyslog struct {
	Facility string `yaml:"facility"`
}
