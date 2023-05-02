package daemon

import (
	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/cmd/localias/shared"
)

var Command = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "daemon",
	Short: "control the proxy server daemon",
	Example: shared.Example(`
# Start the proxy server as a daemon process
localias daemon start
# Show the status of the daemon process
localias daemon status
# Apply the latest configuration to the proxy server in the daemon process
localias daemon reload
# Stop the daemon process
localias daemon stop
	`),
}
