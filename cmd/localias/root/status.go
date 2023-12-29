package root

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/pkg/daemon"
)

var statusCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:     "status",
	Aliases: []string{"ps"},
	Short:   "show the status of the daemon process",
	RunE:    statusImpl,
}

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

func init() { //nolint:gochecknoinits
	Command.AddCommand(statusCmd)
}
