package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func listImpl(_ *cobra.Command, _ []string) error {
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

var listCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "list all managed entries",
	RunE:    listImpl,
}

func init() { //nolint:gochecknoinits
	hostctlCmd.AddCommand(listCmd)
}
