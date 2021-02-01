package commands

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tzvetkoff-go/errors"
)

// NewCompletionCommand ...
func NewCompletionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "completion",
		Short:  "Generates shell completion scripts",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			shell := "bash"
			if len(args) > 0 {
				shell = args[0]
			}

			switch shell {
			case "bash":
				return cmd.Root().GenBashCompletion(os.Stdout)
			case "fish":
				return cmd.Root().GenFishCompletion(os.Stdout, true)
			case "zsh":
				return cmd.Root().GenZshCompletion(os.Stdout)
			case "powershell":
				return cmd.Root().GenPowerShellCompletion(os.Stdout)
			default:
				return errors.New("unknown shell type: %s", shell)
			}
		},
	}

	cmd.Flags().BoolP("help", "h", false, "help message")
	_ = cmd.Flags().MarkHidden("help")

	return cmd
}
