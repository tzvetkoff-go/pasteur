package commands

import (
	"io/fs"

	"github.com/spf13/cobra"

	configPkg "github.com/tzvetkoff-go/pasteur/pkg/config"
	"github.com/tzvetkoff-go/pasteur/pkg/fsutil"
	"github.com/tzvetkoff-go/pasteur/pkg/webserver"
)

// NewExtractStaticCommand ...
func NewExtractStaticCommand() *cobra.Command {
	configPath := configPkg.DefaultConfigPath

	cmd := &cobra.Command{
		Use:   "extract-static <path>",
		Short: "Extract static assets to a directory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			staticFS, err := fs.Sub(webserver.StaticFSRoot, "static")
			if err != nil {
				return err
			}

			return fsutil.Extract(staticFS, args[0])
		},
	}

	cmd.Flags().BoolP("help", "h", false, "help message")
	_ = cmd.Flags().MarkHidden("help")

	cmd.Flags().StringVarP(&configPath, "config", "c", configPath, "Path to config.yml")

	return cmd
}
