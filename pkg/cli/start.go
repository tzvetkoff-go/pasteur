package cli

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tzvetkoff-go/inject"
	"github.com/tzvetkoff-go/logger"

	appPkg "github.com/tzvetkoff-go/pasteur/pkg/app"
	dbPkg "github.com/tzvetkoff-go/pasteur/pkg/db"
	hasherPkg "github.com/tzvetkoff-go/pasteur/pkg/hasher"
	loggerPkg "github.com/tzvetkoff-go/pasteur/pkg/logger"
	webserverPkg "github.com/tzvetkoff-go/pasteur/pkg/webserver"
)

// NewStartCommand ...
func NewStartCommand() *cobra.Command {
	configPath := appPkg.DefaultConfigPath

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the server",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Config ...
			config, err := appPkg.ConfigFromFile(configPath)
			if err != nil {
				logger.Error("%s", err)
				return err
			}

			err = loggerPkg.Configure(&config.Logger)
			if err != nil {
				logger.ClearBackends()
				logger.AddBackends(&logger.WriterBackend{
					Writer:    os.Stderr,
					Formatter: logger.DefaultFormatter,
				})
				logger.Error("%s", err)
			}

			// Injector ...
			injector := inject.New()

			// DB ...
			db, err := dbPkg.New(&config.DB)
			if err != nil {
				logger.Error("%s", err)
				return err
			}
			injector.ProvideObject("DB", db)

			// Hasher ...
			hasher, err := hasherPkg.New(&config.Hasher)
			if err != nil {
				logger.Error("%s", err)
				return err
			}
			injector.ProvideObject("Hasher", hasher)

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

			// Serve ...
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
