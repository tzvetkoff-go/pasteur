package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tzvetkoff-go/pasteur/version"
)

// NewVersionCommand ...
func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		RunE: func(command *cobra.Command, args []string) error {
			fmt.Println(version.Version)
			return nil
		},
	}

	cmd.Flags().BoolP("help", "h", false, "help message")
	_ = cmd.Flags().MarkHidden("help")

	return cmd
}
