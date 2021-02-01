package commands

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tzvetkoff-go/inject"
	"github.com/tzvetkoff-go/logger"

	configPkg "github.com/tzvetkoff-go/pasteur/pkg/config"
	dbPkg "github.com/tzvetkoff-go/pasteur/pkg/db"
	hasherPkg "github.com/tzvetkoff-go/pasteur/pkg/hasher"
	webserverPkg "github.com/tzvetkoff-go/pasteur/pkg/webserver"
)

// NewStartCommand ...
func NewStartCommand() *cobra.Command {
	configPath := configPkg.DefaultConfigPath

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the server",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Config ...
			config, err := configPkg.LoadConfig(configPath)
			if err != nil {
				logger.Error("%s", err)
				return err
			}

			// Logger ...
			logger.SetLevel(logger.LookupString(config.Logger.Level))
			// logger.SetColor(false)

			if config.Logger.Output == "stdout" {
				logger.SetWriter(os.Stdout)
			} else if config.Logger.Output == "stderr" {
				logger.SetWriter(os.Stderr)
			} else if config.Logger.Output == "syslog" {
				logger.Warning("TODO: Implement SysLog logging ...")
			} else {
				file, err := os.OpenFile(config.Logger.Output, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
				if err != nil {
					logger.SetWriter(os.Stderr)
					logger.Warning("could not open log file, falling back to stderr", logger.Fields{
						"logfile": config.Logger.Output,
					})
				}
				logger.SetWriter(file)
			}

			// Injector ...
			injector := inject.New()

			// Hasher ...
			hasher, err := hasherPkg.New(&config.Hasher)
			if err != nil {
				logger.Error("%s", err)
				return err
			}
			injector.ProvideObject("Hasher", hasher)

			// DB ...
			db, err := dbPkg.New(&config.DB)
			if err != nil {
				logger.Error("%s", err)
				return err
			}
			injector.ProvideObject("DB", db)

			// WebServer ...
			webserver, err := webserverPkg.New(&config.WebServer)
			if err != nil {
				logger.Error("%s", err)
				return err
			}
			injector.ProvideObject("WebServer", webserver)

			// Inject ...
			err = injector.Inject(hasher, db, webserver)
			if err != nil {
				logger.Error("%s", err)
				return err
			}

			// Server ...
			err = webserver.Serve()
			if err != nil {
				logger.Error("%s", err)
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolP("help", "h", false, "help message")
	_ = cmd.Flags().MarkHidden("help")

	cmd.Flags().StringVarP(&configPath, "config", "c", configPath, "Path to config.yml")

	return cmd
}
