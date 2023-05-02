package daemon

import (
	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/cmd/localias/shared"
	"github.com/peterldowns/localias/pkg/daemon"
)

func stopImpl(_ *cobra.Command, _ []string) error {
	cfg := shared.Config()
	return daemon.Stop(cfg)
}

var stopCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "stop",
	Short: "stop the daemon process",
	RunE:  stopImpl,
}

func init() { //nolint:gochecknoinits
	Command.AddCommand(stopCmd)
}
