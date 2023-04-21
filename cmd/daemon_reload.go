package cmd

import (
	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/pkg/config"
	"github.com/peterldowns/localias/pkg/daemon"
)

func reloadImpl(_ *cobra.Command, _ []string) error {
	hctl := hostctlController()
	cfg, err := config.Load(nil)
	if err != nil {
		return err
	}
	return daemon.Reload(hctl, cfg)
}

var reloadCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "reload",
	Short: "reload the background daemon's config",
	RunE:  reloadImpl,
}

func init() { //nolint:gochecknoinits
	daemonCmd.AddCommand(reloadCmd)
}
