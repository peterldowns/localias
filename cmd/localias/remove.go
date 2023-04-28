package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func removeImpl(_ *cobra.Command, aliases []string) error {
	cfg := loadConfig()
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

var removeCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:     "remove alias [...more aliases]",
	Aliases: []string{"rm", "delete"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "remove an alias",
	RunE:    removeImpl,
}

func init() { //nolint:gochecknoinits
	rootCmd.AddCommand(removeCmd)
}
