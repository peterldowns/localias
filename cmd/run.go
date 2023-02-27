package cmd

import (
	"github.com/spf13/cobra"

	"github.com/peterldowns/pfpro/pkg/pfpro"
)

func runImpl(_ *cobra.Command, _ []string) error {
	hctl := controller()
	cfg, err := pfpro.Load(nil)
	if err != nil {
		return err
	}
	return pfpro.Run(hctl, cfg)
}

var runCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "run",
	Short: "run the caddy server",
	RunE:  runImpl,
}

func init() { //nolint:gochecknoinits
	rootCmd.AddCommand(runCmd)
}
