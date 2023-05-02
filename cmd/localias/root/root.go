package root

import (
	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/cmd/localias/daemon"
	"github.com/peterldowns/localias/cmd/localias/debug"
	"github.com/peterldowns/localias/cmd/localias/shared"
)

var Command = &cobra.Command{ //nolint:gochecknoglobals
	Version: shared.VersionString(),
	Use:     "localias",
	Short:   "securely proxy domains to local development servers",
	Example: shared.Example(`
# Add an alias forwarding https://secure.test to http://127.0.0.1:9000
localias set secure.test 9000
# Update an existing alias to forward to a different port
localias set secure.test 9001
# Remove an alias
localias remove secure.test
# Show aliases
localias list
# Clear all aliases
localias clear
# Run the server, automatically applying all necessary rules to
# /etc/hosts and creating any necessary TLS certificates
localias run
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
	Command.CompletionOptions.HiddenDefaultCmd = true
	Command.TraverseChildren = true
	Command.SilenceErrors = true
	Command.SilenceUsage = true
	Command.SetVersionTemplate("{{.Version}}\n")

	shared.Flags.Configfile = Command.PersistentFlags().StringP("configfile", "c", "", "path to the configuration file to edit")

	Command.AddCommand(daemon.Command)
	Command.AddCommand(debug.Command)
}
