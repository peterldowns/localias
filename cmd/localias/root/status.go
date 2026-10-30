package root

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/pkg/daemon"
)

var statusCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:     "status",
	Aliases: []string{"ps"},
	Short:   "show the status of the daemon process",
	Long:    "Show whether the background daemon process is running. This does not check proxy readiness or whether saved configuration has been applied.",
	Args:    cobra.NoArgs,
	RunE:    statusImpl,
}

func statusImpl(cmd *cobra.Command, _ []string) error {
	jsonOutput, err := cmd.Flags().GetBool("json")
	if err != nil {
		return err
	}
	proc, err := daemon.Status()
	if err != nil {
		return err
	}
	var pid *int
	if proc != nil {
		pid = &proc.Pid
		defer func() { _ = proc.Release() }()
	}
	w := cmd.OutOrStdout()
	if jsonOutput {
		output := struct {
			Running bool `json:"running"`
			PID     *int `json:"pid"`
		}{Running: proc != nil, PID: pid}
		return json.NewEncoder(w).Encode(output)
	}
	if proc == nil {
		_, err = fmt.Fprintln(w, "daemon is not running")
	} else {
		_, err = fmt.Fprintf(w, "daemon running with pid %d\n", proc.Pid)
	}
	return err
}

func init() {
	Command.AddCommand(statusCmd)
	statusCmd.Flags().Bool("json", false, "output daemon process status as JSON")
}
