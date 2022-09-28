package app

import (
	"io"
	"os"

	"github.com/tzvetkoff-go/errors"
	"gopkg.in/yaml.v2"
)

// ConfigFromFile ...
func ConfigFromFile(path string) (*Config, error) {
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
