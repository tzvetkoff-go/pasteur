package commands

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/tzvetkoff-go/pasteur/version"
)

// NewRootCommand ...
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "pasteur",
		Short:        "pasteur",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd.Help()
			os.Exit(1) // revive:disable-line:deep-exit
			return nil
		},
	}

	cmd.Flags().BoolP("help", "h", false, "help message")
	_ = cmd.Flags().MarkHidden("help")

	cmd.SetVersionTemplate(version.Version)
	cmd.AddCommand(
		NewStartCommand(),
		NewHashEncodeCommand(),
		NewHashDecodeCommand(),
		NewExtractStaticCommand(),
		NewExtractTemplatesCommand(),
		NewVersionCommand(),
		NewCompletionCommand(),
	)

	return cmd
}
