package logger

import (
	"os"
	"path"

	"github.com/tzvetkoff-go/logger"
	"github.com/tzvetkoff-go/logger/backends/syslog"
)

// Configure ...
func Configure(config *Config) error {
	logger.ClearBackends()
	logger.SetLevel(logger.LookupString(config.Level))
	// logger.SetColor(false)

	for _, backend := range config.Backends {
		switch backend.Type {
		case "file":
			var file *os.File
			var err error

			switch backend.File.Destination {
			case "stdout":
				file, err = os.Stdout, nil
			case "stderr":
				file, err = os.Stderr, nil
			default:
				file, err = os.OpenFile(backend.File.Destination, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			}

			if err != nil {
				return err
			}

			logger.AddBackends(&logger.WriterBackend{
				Writer:    file,
				Formatter: logger.DefaultFormatter,
			})
		case "syslog":
			priority := syslog.Priority(logger.DefaultLogger.Level) | syslog.LookupString(backend.Syslog.Facility)
			backend, err := syslog.NewBackend(
				"",
				"",
				priority,
				path.Base(os.Args[0]),
				&logger.JSONFormatter{},
			)

			if err != nil {
				return err
			}

			logger.AddBackends(backend)
		}
	}

	return nil
}
