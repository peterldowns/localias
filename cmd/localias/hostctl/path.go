package hostctl

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/cmd/localias/shared"
	"github.com/peterldowns/localias/pkg/hostctl"
)

func pathImpl(_ *cobra.Command, _ []string) error {
	hctl := shared.Controller()
	for _, path := range getPaths(hctl) {
		fmt.Println(path)
	}
	return nil
}

var pathCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:     "path",
	Aliases: []string{"paths"},
	Short:   "show the path(s) of the host file(s) being edited by localias",
	RunE:    pathImpl,
}

func init() { //nolint:gochecknoinits
	Command.AddCommand(pathCmd)
}

func getPaths(hctl hostctl.Controller) []string {
	switch c := hctl.(type) {
	case *hostctl.FileController:
		return []string{c.Path}
	case *hostctl.WindowsController:
		return []string{c.WindowsHostsFile}
	case *hostctl.WSLController:
		return []string{
			// TODO: extract to a constant and a helper for getting the unix
			// version of this
			`$env:windir\System32\drivers\etc\hosts`,
		}
	case hostctl.MultiController:
		var paths []string
		for _, ctl := range c {
			paths = append(paths, getPaths(ctl)...)
		}
		return paths
	default:
		return []string{}
	}
}
