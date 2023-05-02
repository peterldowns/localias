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
# print the contents of the current configuration file
localias debug config --print
# show the path to the root certificate
localias debug cert
# print the contents of the root certificate
localias debug cert --print

	`),
	Hidden: true,
}

func init() { //nolint:gochecknoinits
	rootCmd.AddCommand(debugCmd)
}
