package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/peterldowns/pfpro/pkg/pfpro"
)

func removeImpl(_ *cobra.Command, aliases []string) error {
	cfg, err := pfpro.Load(nil)
	if err != nil {
		return err
	}
	removed := []pfpro.Directive{}
	preserved := []pfpro.Directive{}
outer:
	for _, d := range cfg.Directives {
		for _, alias := range aliases {
			if d.Upstream == alias {
				removed = append(removed, d)
				continue outer
			}
		}
		preserved = append(preserved, d)
	}
	cfg.Directives = preserved
	if err := pfpro.WriteConfig(cfg); err != nil {
		return err
	}
	for _, d := range removed {
		fmt.Printf(
			"%s %s -> %s\n",
			color.New(color.FgRed).Sprint("[removed]"),
			color.New(color.FgBlue).Sprint(d.Upstream),
			color.New(color.FgWhite).Sprint(d.Downstream),
		)
	}
	return nil
}

var removeCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:     "remove [aliases...]",
	Aliases: []string{"rm", "delete"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "remove aliases",
	RunE:    removeImpl,
}

func init() { //nolint:gochecknoinits
	rootCmd.AddCommand(removeCmd)
}
