package root

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/cmd/localias/shared"
)

var listCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:     "list",
	Aliases: []string{"l", "ls", "show"},
	Short:   "list all aliases",
	RunE:    listImpl,
}

func listImpl(_ *cobra.Command, _ []string) error {
	cfg := shared.Config()
	for _, entry := range cfg.Entries {
		fmt.Printf(
			"%s -> %s\n",
			color.New(color.FgBlue).Sprint(entry.Alias),
			color.New(color.FgWhite).Sprintf("%d", entry.Port),
		)
	}
	return nil
}

func init() { //nolint:gochecknoinits
	Command.AddCommand(listCmd)
}
