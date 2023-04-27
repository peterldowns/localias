package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func configImpl(_ *cobra.Command, _ []string) error {
	cfg := loadConfig()
	fmt.Println(cfg.Path)
	return nil
}

var configCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "config",
	Short: "show the configuration file path",
	RunE:  configImpl,
}

func init() { //nolint:gochecknoinits
	rootCmd.AddCommand(configCmd)
}
