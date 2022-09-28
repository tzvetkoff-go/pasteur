package cli

import (
	"io/fs"

	"github.com/spf13/cobra"

	"github.com/tzvetkoff-go/pasteur/pkg/fsutil"
	"github.com/tzvetkoff-go/pasteur/pkg/webserver"
)

// NewExtractTemplatesCommand ...
func NewExtractTemplatesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extract-templates <path>",
		Short: "Extract templates to a directory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			templatesFS, err := fs.Sub(webserver.TemplatesFSRoot, "templates")
			if err != nil {
				return err
			}

			return fsutil.Extract(templatesFS, args[0])
		},
	}

	cmd.Flags().BoolP("help", "h", false, "help message")
	_ = cmd.Flags().MarkHidden("help")

	return cmd
}
