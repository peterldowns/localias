package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func debugImpl(_ *cobra.Command, _ []string) error {
	cfg := loadConfig()
	fmt.Println("--- configfile:")
	fmt.Println(cfg.Path)

	fmt.Println("--- config entries:")
	if err := listImpl(nil, nil); err != nil {
		return err
	}

	fmt.Println("--- hostctl entries")
	if err := hostctlListImpl(nil, nil); err != nil {
		return err
	}
	fmt.Println("--- caddyfile")
	fmt.Println(cfg.Caddyfile())
	return nil
}

var debugCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:    "debug",
	Short:  "debug the configuration",
	RunE:   debugImpl,
	Hidden: true,
}

func init() { //nolint:gochecknoinits
	rootCmd.AddCommand(debugCmd)
}
