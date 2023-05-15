package hostctl

import (
	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/cmd/localias/shared"
	"github.com/peterldowns/localias/pkg/config"
)

func applyImpl(_ *cobra.Command, _ []string) error {
	cfg := shared.Config()
	hctl := shared.Controller()
	if err := config.Apply(hctl, cfg); err != nil {
		return err
	}
	return nil
}

var applyCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:     "apply",
	Aliases: []string{"sync", "update"},
	Short:   "apply the current configuration to the hosts file(s)",
	RunE:    applyImpl,
}

func init() { //nolint:gochecknoinits
	Command.AddCommand(applyCmd)
}
