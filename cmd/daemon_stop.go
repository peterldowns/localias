package cmd

import (
	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/pkg/config"
	"github.com/peterldowns/localias/pkg/daemon"
)

func stopImpl(_ *cobra.Command, _ []string) error {
	hctl := hostctlController()
	cfg, err := config.Load(nil)
	if err != nil {
		return err
	}
	return daemon.Stop(hctl, cfg)
}

var stopCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "stop",
	Short: "stop running the background daemon",
	RunE:  stopImpl,
}

func init() { //nolint:gochecknoinits
	daemonCmd.AddCommand(stopCmd)
}
