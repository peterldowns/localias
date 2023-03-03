package cmd

import (
	"github.com/spf13/cobra"

	"github.com/peterldowns/pfpro/pkg/config"
	"github.com/peterldowns/pfpro/pkg/server"
)

func runImpl(_ *cobra.Command, _ []string) error {
	hctl := hostctlController()
	cfg, err := config.Load(nil)
	if err != nil {
		return err
	}
	return server.Run(hctl, cfg)
}

var runCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "run",
	Short: "run the caddy server",
	RunE:  runImpl,
}

func init() { //nolint:gochecknoinits
	rootCmd.AddCommand(runCmd)
}
