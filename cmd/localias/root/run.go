package root

import (
	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/cmd/localias/shared"
	"github.com/peterldowns/localias/pkg/config"
	"github.com/peterldowns/localias/pkg/server"
)

func runImpl(_ *cobra.Command, _ []string) error {
	hctl := shared.Controller()
	cfg := shared.Config()
	if err := config.Apply(hctl, cfg); err != nil {
		return err
	}
	if err := server.Start(cfg); err != nil {
		return err
	}
	select {}
}

var runCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "run",
	Short: "run the proxy server in the foreground",
	RunE:  runImpl,
}

func init() { //nolint:gochecknoinits
	Command.AddCommand(runCmd)
}
