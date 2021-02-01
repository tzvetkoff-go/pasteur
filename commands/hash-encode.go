package commands

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	configPkg "github.com/tzvetkoff-go/pasteur/pkg/config"
	hasherPkg "github.com/tzvetkoff-go/pasteur/pkg/hasher"
)

// NewHashEncodeCommand ...
func NewHashEncodeCommand() *cobra.Command {
	configPath := configPkg.DefaultConfigPath

	cmd := &cobra.Command{
		Use:   "hash-encode <id>",
		Short: "Encode integer to hash",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Config ...
			config, err := configPkg.LoadConfig(configPath)
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
