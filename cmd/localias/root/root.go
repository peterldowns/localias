package root

import (
	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/cmd/localias/daemon"
	"github.com/peterldowns/localias/cmd/localias/debug"
	"github.com/peterldowns/localias/cmd/localias/hostctl"
	"github.com/peterldowns/localias/cmd/localias/shared"
)

var Command = &cobra.Command{ //nolint:gochecknoglobals
	Version: shared.VersionString(),
	Use:     "localias",
	Short:   "securely manage local aliases for development servers",
	Example: shared.Example(`
# Add an alias forwarding https://secure.test to http://127.0.0.1:9000
localias set secure.test 9000
# Update an existing alias to forward to a different port
localias set secure.test 9001
# Remove an alias
localias remove secure.test
# List all aliases
localias list
# Clear all aliases
localias clear

# Run the proxy server in the foreground
localias run
# Start the proxy server as a daemon process
localias daemon start
# Show the status of the daemon process
localias daemon status
# Apply the latest configuration to the proxy server in the daemon process
localias daemon reload
# Stop the daemon process
localias daemon stop

# Show the host file(s) that localias edits
localias hostctl print
# Show the entries that localias has added to the host file(s)
localias hostctl list
# Remove all localias-managed entries from the host file(s)
localias hostctl clear
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
	Command.AddCommand(hostctl.Command)
}
