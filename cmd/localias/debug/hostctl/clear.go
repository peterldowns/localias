package hostctl

import (
	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/cmd/localias/shared"
)

func clearImpl(_ *cobra.Command, _ []string) error {
	hctl := shared.Controller()
	if err := hctl.Clear(); err != nil {
		return err
	}
	if _, err := hctl.Apply(); err != nil {
		return err
	}
	return nil
}

var clearCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:     "clear",
	Aliases: []string{"delete", "remove"},
	Short:   "clear all localias-managed entries from the host file(s)",
	RunE:    clearImpl,
}

func init() {
	Command.AddCommand(clearCmd)
}
