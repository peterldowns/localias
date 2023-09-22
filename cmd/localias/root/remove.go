package root

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/cmd/localias/shared"
)

var removeCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:     "rm alias [...more aliases]",
	Aliases: []string{"remove", "delete"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "remove an alias",
	RunE:    removeImpl,
}

func removeImpl(_ *cobra.Command, aliases []string) error {
	cfg := shared.Config()
	removed := cfg.Remove(aliases...)
	if err := cfg.Save(); err != nil {
		return err
	}
	for _, d := range removed {
		fmt.Printf(
			"%s %s -> %s\n",
			color.New(color.FgRed).Sprint("[removed]"),
			color.New(color.FgBlue).Sprint(d.Alias),
			color.New(color.FgWhite).Sprint(d.Port),
		)
	}
	return nil
}

func init() { //nolint:gochecknoinits
	Command.AddCommand(removeCmd)
}
