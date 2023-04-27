package cmd

import (
	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/pkg/daemon"
	"github.com/peterldowns/localias/pkg/hostctl"
)

func runImpl(_ *cobra.Command, _ []string) error {
	hctl := hostctl.DefaultController()
	cfg := loadConfig()
	return daemon.Run(hctl, cfg)
}

var runCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "run",
	Short: "run the caddy server",
	RunE:  runImpl,
}

func init() { //nolint:gochecknoinits
	rootCmd.AddCommand(runCmd)
}
