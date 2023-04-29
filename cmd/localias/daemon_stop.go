package main

import (
	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/pkg/daemon"
)

func stopImpl(_ *cobra.Command, _ []string) error {
	cfg := loadConfig()
	return daemon.Stop(cfg)
}

var stopCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "stop",
	Short: "stop running the background daemon",
	RunE:  stopImpl,
}

func init() { //nolint:gochecknoinits
	daemonCmd.AddCommand(stopCmd)
}
