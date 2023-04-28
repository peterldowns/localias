package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/pkg/daemon"
)

func statusImpl(_ *cobra.Command, _ []string) error {
	proc, err := daemon.Status()
	if err != nil {
		return err
	}
	if proc == nil {
		fmt.Printf("daemon is not running\n")
	} else {
		fmt.Printf("daemon running with pid %d\n", proc.Pid)
	}
	return nil
}

var statusCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "status",
	Short: "show the status of the background daemon",
	RunE:  statusImpl,
}

func init() { //nolint:gochecknoinits
	daemonCmd.AddCommand(statusCmd)
}
