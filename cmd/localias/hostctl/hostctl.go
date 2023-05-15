package hostctl

import (
	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/cmd/localias/shared"
)

var Command = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "hostctl",
	Short: "interact with the hosts file(s) that localias manages",
	Example: shared.Example(`
# Show the path(s) of the host file(s)
localias hostctl path
# Show any entries that localias has made to the host file(s)
localias hostctl list
# Apply the current configuration to the host file(s)
localias hostctl apply
# Clear all entries from the host file(s)
localias hostctl clear
	`),
}
