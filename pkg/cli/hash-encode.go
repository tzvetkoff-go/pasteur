package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	appPkg "github.com/tzvetkoff-go/pasteur/pkg/app"
	hasherPkg "github.com/tzvetkoff-go/pasteur/pkg/hasher"
)

// NewHashEncodeCommand ...
func NewHashEncodeCommand() *cobra.Command {
	configPath := appPkg.DefaultConfigPath

	cmd := &cobra.Command{
		Use:   "hash-encode <id>",
		Short: "Encode integer to hash",
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

			// Parse arg ...
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			// Encode ...
			hash, err := hasher.Encode(id)
			if err != nil {
				return err
			}

			// Print
			fmt.Println(hash)

			return nil
		},
	}

	cmd.Flags().BoolP("help", "h", false, "help message")
	_ = cmd.Flags().MarkHidden("help")

	cmd.Flags().StringVarP(&configPath, "config", "c", configPath, "Path to config.yml")

	return cmd
}
