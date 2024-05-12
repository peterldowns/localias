package root

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/pkg/daemon"
)

var stopCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "stop",
	Short: "stop the daemon process",
	RunE:  stopImpl,
}

func stopImpl(_ *cobra.Command, _ []string) error {
	existing, err := daemon.Status()
	if err != nil {
		return err
	}
	if existing == nil {
		fmt.Println("daemon is not running")
		return nil
	}

	fmt.Printf("stopping daemon on pid %d\n", existing.Pid)
	return existing.Kill()
}

func init() {
	Command.AddCommand(stopCmd)
}
