package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	appPkg "github.com/tzvetkoff-go/pasteur/pkg/app"
	hasherPkg "github.com/tzvetkoff-go/pasteur/pkg/hasher"
)

// NewHashDecodeCommand ...
func NewHashDecodeCommand() *cobra.Command {
	configPath := appPkg.DefaultConfigPath

	cmd := &cobra.Command{
		Use:   "hash-decode <hash>",
		Short: "Decode hash to integer",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Config ...
			config, err := appPkg.ConfigFromFile(configPath)
			if err != nil {
				return err
			}

			// Hasher ...
			hasher, err := hasherPkg.New(&config.Hasher)
			if err != nil {
				return err
			}

			// Decode ...
			id, err := hasher.Decode(args[0])
			if err != nil {
				return err
			}

			// Print
			fmt.Println(id)

			return nil
		},
	}

	cmd.Flags().BoolP("help", "h", false, "help message")
	_ = cmd.Flags().MarkHidden("help")

	cmd.Flags().StringVarP(&configPath, "config", "c", configPath, "Path to config.yml")

	return cmd
}
