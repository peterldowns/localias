package cmd

import (
	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/pkg/config"
	"github.com/peterldowns/localias/pkg/daemon"
)

func startImpl(_ *cobra.Command, _ []string) error {
	hctl := hostctlController()
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
