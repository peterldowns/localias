package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/pkg/config"
)

func listImpl(_ *cobra.Command, _ []string) error {
	cfg, err := config.Load(nil)
	if err != nil {
		return err
	}
	for _, directive := range cfg.Directives {
		fmt.Printf(
			"%s -> %s\n",
			color.New(color.FgBlue).Sprint(directive.Upstream),
			color.New(color.FgWhite).Sprint(directive.Downstream),
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
