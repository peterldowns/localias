package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addFlags struct { //nolint:gochecknoglobals
	Force *bool
}

func addImpl(_ *cobra.Command, args []string) error {
	c := controller()

	var ip string
	var aliases []string
	if len(args) == 1 {
		ip = "127.0.0.1"
		aliases = args
	} else {
		ip = args[0]
		aliases = args[1:]
	}
	lines, err := c.Add(*addFlags.Force, ip, aliases...)
	if err != nil {
		return err
	}
	if err := c.Save(); err != nil {
		return err
	}
	for _, line := range lines {
		fmt.Println("[added] ", line)
	}
	return nil
}

var addCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:     "add [IP address] [aliases...]",
	Aliases: []string{"a", "new", "create"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "add a new managed entry",
	RunE:    addImpl,
}

func init() { //nolint:gochecknoinits
	addFlags.Force = addCmd.Flags().BoolP("force", "f", false, "on conflict, remove existing rules")
	hostctlCmd.AddCommand(addCmd)
}
