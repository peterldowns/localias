package root

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/cmd/localias/shared"
	"github.com/peterldowns/localias/pkg/config"
	"github.com/peterldowns/localias/pkg/daemon"
	"github.com/peterldowns/localias/pkg/server"
)

var runCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "run",
	Short: "run the proxy server in the foreground",
	RunE:  runImpl,
}

func runImpl(_ *cobra.Command, _ []string) error {
	// Apply the config to /etc/hosts
	hctl := shared.Controller()
	cfg := shared.Config()
	if err := config.Apply(hctl, cfg); err != nil {
		return err
	}
	// If the daemon is already running, print a warning and then kill it so
	// that this instance can run instead.
	existing, err := daemon.Status()
	if err != nil {
		return err
	}
	if existing != nil {
		fmt.Printf("replacing existing daemon on pid %d\n", existing.Pid)
		if err := daemon.Kill(); err != nil {
			return err
		}
	}
	// Start the servers.
	instance := &server.Server{Config: cfg}
	if err := instance.Start(); err != nil {
		return err
	}
	// Loop and run until sigint/sigterm/sigabrt
	server.WaitForExitSignal()
	return nil
}

func init() {
	Command.AddCommand(runCmd)
}
