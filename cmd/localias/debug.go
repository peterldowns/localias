package main

import (
	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/pkg/util"
)

var debugCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "debug",
	Short: "various helpers for debugging localias",
	Example: util.Example(`
# show the path to the current configuration file
localias debug config
	`),
}

func init() { //nolint:gochecknoinits
	rootCmd.AddCommand(debugCmd)
}
