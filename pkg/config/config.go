package config

import (
	"io"
	"os"

	"github.com/tzvetkoff-go/errors"
	"gopkg.in/yaml.v2"
)

// DefaultConfigPath ...
const DefaultConfigPath = `./config.yml`

// Config ...
type Config struct {
	Logger    Logger    `yaml:"logger"`
	DB        DB        `yaml:"db"`
	Hasher    Hasher    `yaml:"hasher"`
	WebServer WebServer `yaml:"webserver"`
}

// Logger ...
type Logger struct {
	Output string `yaml:"output"`
	Level  string `yaml:"level"`
}

// DB ...
type DB struct {
	Storage    string     `yaml:"storage"`
	FileSystem FileSystem `yaml:"filesystem"`
	SQL        SQL        `yaml:"sql"`
}

// FileSystem ...
type FileSystem struct {
	Root string `yaml:"root"`
}

// SQL ...
type SQL struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}

// Hasher ...
type Hasher struct {
	Salt      string `yaml:"salt"`
	Alphabet  string `yaml:"alphabet"`
	MinLength int    `yaml:"min-length"`
}

// WebServer ...
type WebServer struct {
	StaticPath      string `yaml:"static-path"`
	TemplatesPath   string `yaml:"templates-path"`
	ListenAddress   string `yaml:"listen-address"`
	ProxyHeader     string `yaml:"proxy-header"`
	TLSCert         string `yaml:"tls-cert"`
	TLSKey          string `yaml:"tls-key"`
	RelativeURLRoot string `yaml:"relative-url-root"`
}

// LoadConfig ...
func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Propagate(err, "could not open config file")
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.Propagate(err, "could not read config file")
	}

	result := &Config{}
	err = yaml.Unmarshal(data, &result)
	if err != nil {
		return nil, errors.Propagate(err, "could not parse config file")
	}

	return result, nil
}
