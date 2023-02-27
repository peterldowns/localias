package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func hostctlListImpl(_ *cobra.Command, _ []string) error {
	c := controller()
	lines, err := c.List()
	if err != nil {
		return err
	}
	for _, line := range lines {
		fmt.Println(line.String())
	}
	return nil
}

var hostctlListCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "list all managed entries",
	RunE:    hostctlListImpl,
}

func init() { //nolint:gochecknoinits
	hostctlCmd.AddCommand(hostctlListCmd)
}
