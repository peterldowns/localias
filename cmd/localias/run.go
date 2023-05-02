package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/pkg/hostctl"
	"github.com/peterldowns/localias/pkg/server"
)

func runImpl(_ *cobra.Command, _ []string) error {
	hctl := hostctl.NewWSL2Controller()
	fmt.Println(hctl.TmpController.HostsFile)
	cfg := loadConfig()
	if err := server.Start(hctl, cfg); err != nil {
		return err
	}
	fmt.Println("applied config")
	select {}
}

var runCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "run",
	Short: "run the caddy server",
	RunE:  runImpl,
}

func init() { //nolint:gochecknoinits
	rootCmd.AddCommand(runCmd)
}
