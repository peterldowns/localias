/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/peterldowns/pfpro/pkg/hostctl"
)

var hostctlFlags struct { //nolint:gochecknoglobals
	File   *string
	Name   *string
	Sudo   *bool
	DryRun *bool
}

func hostctlController() *hostctl.Controller {
	return hostctl.NewController(
		*hostctlFlags.File,
		*hostctlFlags.Sudo,
		*hostctlFlags.DryRun,
		*hostctlFlags.Name,
	)
}

var hostctlCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:    "hostctl",
	Hidden: true,
	Short:  "modify an /etc/hosts-type file",
}

var hostctlAddCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:     "add alias [ip address]",
	Aliases: []string{"a", "new", "create"},
	Args:    cobra.RangeArgs(1, 2),
	Short:   "add a new entry",
	RunE:    hostctlAddImpl,
}

var hostctlListCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "list all entries",
	RunE:    hostctlListImpl,
}

var hostctlRemoveCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:     "remove [aliases...]",
	Aliases: []string{"a", "new", "create"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "remove entries",
	RunE:    hostctlRemoveImpl,
}

func init() { //nolint:gochecknoinits
	rootCmd.AddCommand(hostctlCmd)
	hostctlFlags.File = hostctlCmd.PersistentFlags().String("file", hostctl.DefaultHostsFile, "which hosts file to modify")
	hostctlFlags.Sudo = hostctlCmd.PersistentFlags().Bool("sudo", hostctl.DefaultSudo, "use sudo to write the hosts file")
	hostctlFlags.Name = hostctlCmd.PersistentFlags().String("name", hostctl.DefaultName, "controller name")
	hostctlFlags.DryRun = hostctlCmd.PersistentFlags().Bool("dry-run", hostctl.DefaultDryRun, "dry run")
	hostctlCmd.AddCommand(hostctlAddCmd)
	hostctlCmd.AddCommand(hostctlListCmd)
	hostctlCmd.AddCommand(hostctlRemoveCmd)
}

func hostctlAddImpl(_ *cobra.Command, args []string) error {
	alias := args[0]
	var ip string
	if len(args) == 2 {
		ip = args[1]
	} else {
		ip = "127.0.0.1"
	}

	c := hostctlController()
	if err := c.Set(ip, alias); err != nil {
		return err
	}
	if err := c.Apply(); err != nil {
		return err
	}
	fmt.Printf("[added] %s -> %s\n", alias, ip)
	return nil
}

func hostctlListImpl(_ *cobra.Command, _ []string) error {
	c := hostctlController()
	lines, err := c.List()
	if err != nil {
		return err
	}
	for _, line := range lines {
		fmt.Println(line.String())
	}
	return nil
}

func hostctlRemoveImpl(_ *cobra.Command, aliases []string) error {
	c := hostctlController()
	for _, alias := range aliases {
		if err := c.Remove(alias); err != nil {
			return err
		}
	}
	if err := c.Apply(); err != nil {
		return err
	}
	for _, alias := range aliases {
		fmt.Printf("[removed] %s\n", alias)
	}
	return nil
}
