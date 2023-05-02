package daemon

import (
	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/cmd/localias/shared"
	"github.com/peterldowns/localias/pkg/config"
	"github.com/peterldowns/localias/pkg/daemon"
)

func reloadImpl(_ *cobra.Command, _ []string) error {
	hctl := shared.Controller()
	cfg := shared.Config()
	if err := config.Apply(hctl, cfg); err != nil {
		return err
	}
	return daemon.Reload(cfg)
}

var reloadCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "reload",
	Short: "apply the latest configuration to the proxy server in the daemon process",
	RunE:  reloadImpl,
}

func init() { //nolint:gochecknoinits
	Command.AddCommand(reloadCmd)
}
