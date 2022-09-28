package cli

import (
	"github.com/spf13/cobra"
	"github.com/tzvetkoff-go/logger"

	appPkg "github.com/tzvetkoff-go/pasteur/pkg/app"
	dbPkg "github.com/tzvetkoff-go/pasteur/pkg/db"
)

// NewMigrateCommand ...
func NewMigrateCommand() *cobra.Command {
	configPath := appPkg.DefaultConfigPath

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate database",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Config ...
			config, err := appPkg.ConfigFromFile(configPath)
			if err != nil {
				return err
			}

			// DB ...
			db, err := dbPkg.New(&config.DB)
			if err != nil {
				logger.Error("%s", err)
				return err
			}

			return db.Migrate()
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}
	cmd.SilenceUsage = false

	cmd.Flags().BoolP("help", "h", false, "help message")
	_ = cmd.Flags().MarkHidden("help")

	cmd.Flags().StringVarP(&configPath, "config", "c", configPath, "Path to config.yml")

	return cmd
}
