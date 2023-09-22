package root

import (
	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/cmd/localias/shared"
	"github.com/peterldowns/localias/pkg/config"
	"github.com/peterldowns/localias/pkg/daemon"
)

var reloadCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:     "reload",
	Aliases: []string{"restart"},
	Short:   "apply the latest configuration to the proxy server in the daemon process",
	RunE:    reloadImpl,
}

func reloadImpl(_ *cobra.Command, _ []string) error {
	// Ensure that the daemon is running.
	existing, err := daemon.Status()
	if err != nil {
		return err
	}
	if existing == nil {
		return shared.DaemonNotRunning{}
	}
	// Apply the config to /etc/hosts
	hctl := shared.Controller()
	cfg := shared.Config()
	if err := config.Apply(hctl, cfg); err != nil {
		return err
	}
	// Reload the daemon with the new config.
	return daemon.Reload(cfg)
}

func init() { //nolint:gochecknoinits
	Command.AddCommand(reloadCmd)
}
