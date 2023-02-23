package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func removeImpl(_ *cobra.Command, aliases []string) error {
	c := controller()
	removed, err := c.Remove(aliases...)
	if err != nil {
		return err
	}
	if err := c.Save(); err != nil {
		return err
	}
	for _, line := range removed {
		fmt.Println("[removed] ", line)
	}
	return nil
}

var removeCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:     "remove [aliases...]",
	Aliases: []string{"a", "new", "create"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "remove a new managed entry",
	RunE:    removeImpl,
}

func init() { //nolint:gochecknoinits
	hostctlCmd.AddCommand(removeCmd)
}
