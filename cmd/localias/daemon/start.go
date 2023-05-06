package daemon

import (
	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/cmd/localias/shared"
	"github.com/peterldowns/localias/pkg/config"
	"github.com/peterldowns/localias/pkg/daemon"
)

func startImpl(_ *cobra.Command, _ []string) error {
	hctl := shared.Controller()
	cfg := shared.Config()
	if err := config.Apply(hctl, cfg); err != nil {
		return err
	}
	return daemon.Start(cfg)
}

var startCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:     "start",
	Aliases: []string{"run", "launch"},
	Short:   "start the proxy server as a daemon process",
	RunE:    startImpl,
}

func init() { //nolint:gochecknoinits
	Command.AddCommand(startCmd)
}
