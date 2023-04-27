package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func listImpl(_ *cobra.Command, _ []string) error {
	cfg := loadConfig()
	for _, entry := range cfg.Entries {
		fmt.Printf(
			"%s -> %s\n",
			color.New(color.FgBlue).Sprint(entry.Alias),
			color.New(color.FgWhite).Sprintf("%d", entry.Port),
		)
	}
	return nil
}

var listCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "list all aliases",
	RunE:    listImpl,
}

func init() { //nolint:gochecknoinits
	rootCmd.AddCommand(listCmd)
}
