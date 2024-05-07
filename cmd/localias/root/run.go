package root

import (
	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/cmd/localias/shared"
	"github.com/peterldowns/localias/pkg/config"
	"github.com/peterldowns/localias/pkg/daemon"
)

var runCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "run",
	Short: "run the proxy server in the foreground",
	RunE:  runImpl,
}

func runImpl(_ *cobra.Command, _ []string) error {
	// Ensure that the daemon is not already running.
	existing, err := daemon.Status()
	if err != nil {
		return err
	}
	if existing != nil {
		return shared.DaemonRunningError{Pid: existing.Pid}
	}
	// Apply the config to /etc/hosts
	hctl := shared.Controller()
	cfg := shared.Config()
	if err := config.Apply(hctl, cfg); err != nil {
		return err
	}
	return daemon.Run(cfg)
}

func init() {
	Command.AddCommand(runCmd)
}
