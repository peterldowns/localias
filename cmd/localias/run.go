package main

import (
	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/pkg/hostctl"
	"github.com/peterldowns/localias/pkg/server"
)

func runImpl(_ *cobra.Command, _ []string) error {
	hctl := hostctl.DefaultController()
	cfg := loadConfig()
	if err := server.Start(hctl, cfg); err != nil {
		return err
	}
	select {} //nolint:revive // valid empty block, keeps the server running forever.
}

var runCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "run",
	Short: "run the caddy server",
	RunE:  runImpl,
}

func init() { //nolint:gochecknoinits
	rootCmd.AddCommand(runCmd)
}
