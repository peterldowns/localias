package main

import (
	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/pkg/util"
)

var daemonCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "daemon",
	Short: "interact with the daemon process",
	Example: util.Example(`
# Run the server as a daemon
localias daemon start
# Check whether or not the daemon is running
localias daemon status
# Reload the config that the daemon is using
localias daemon reload
# Stop the daemon if it is running
localias daemon stop
	`),
}

func init() { //nolint:gochecknoinits
	rootCmd.AddCommand(daemonCmd)
}
