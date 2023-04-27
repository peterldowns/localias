package cmd

import (
	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/pkg/config"
	"github.com/peterldowns/localias/pkg/daemon"
	"github.com/peterldowns/localias/pkg/hostctl"
)

func startImpl(_ *cobra.Command, _ []string) error {
	hctl := hostctl.DefaultController()
	cfg, err := config.Load(nil)
	if err != nil {
		return err
	}
	return daemon.Start(hctl, cfg)
}

var startCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "start",
	Short: "start running the background daemon",
	RunE:  startImpl,
}

func init() { //nolint:gochecknoinits
	daemonCmd.AddCommand(startCmd)
}
