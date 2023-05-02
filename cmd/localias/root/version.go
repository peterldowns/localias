package root

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/cmd/localias/shared"
)

func versionImpl(_ *cobra.Command, _ []string) error {
	fmt.Println(shared.VersionString())
	return nil
}

var versionCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "version",
	Short: "show the version of this binary",
	RunE:  versionImpl,
}

func init() { //nolint:gochecknoinits
	Command.AddCommand(versionCmd)
}
