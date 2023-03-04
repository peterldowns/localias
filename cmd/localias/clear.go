package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/pkg/config"
)

func clearImpl(_ *cobra.Command, _ []string) error {
	cfg, err := config.Load(nil)
	if err != nil {
		return err
	}
	removed := cfg.Directives
	cfg.Directives = nil
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

var clearCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "clear",
	Short: "clear all aliases",
	RunE:  clearImpl,
}

func init() { //nolint:gochecknoinits
	rootCmd.AddCommand(clearCmd)
}
