package root

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/cmd/localias/debug"
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
localias rm secure.test
# List all aliases
localias list
# Clear all aliases
localias clear

# Start the proxy server as a daemon process
localias start
# Show the status of the daemon process
localias status
# Apply the latest configuration and relaunch the daemon process
localias reload
# Stop the daemon process
localias stop
# Run the proxy server in the foreground
localias run
  `),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return fmt.Errorf(`invalid command: "%s"`, args[0])
		}
		return cmd.Help()
	},
}

func init() {
	Command.CompletionOptions.HiddenDefaultCmd = true
	Command.TraverseChildren = true
	Command.SilenceErrors = true
	Command.SilenceUsage = true
	Command.SetVersionTemplate("{{.Version}}\n")

	shared.Flags.Configfile = Command.PersistentFlags().StringP("configfile", "c", "", "path to the configuration file to edit")

	Command.AddCommand(debug.Command)
}
