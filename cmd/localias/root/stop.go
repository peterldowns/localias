package root

import (
	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/cmd/localias/shared"
	"github.com/peterldowns/localias/pkg/daemon"
)

var stopCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "stop",
	Short: "stop the daemon process",
	RunE:  stopImpl,
}

func stopImpl(_ *cobra.Command, _ []string) error {
	// Ensure that the daemon is running .
	existing, err := daemon.Status()
	if err != nil {
		return err
	}
	if existing == nil {
		return shared.DaemonNotRunning{}
	}
	return existing.Kill()
}

func init() { //nolint:gochecknoinits
	Command.AddCommand(stopCmd)
}
