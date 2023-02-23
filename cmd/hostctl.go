/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/peterldowns/pfpro/pkg/hostctl"
)

var hostctlFlags struct { //nolint:gochecknoglobals
	File   *string
	Name   *string
	Sudo   *bool
	DryRun *bool
}

func controller() *hostctl.Controller {
	return hostctl.NewController(
		*hostctlFlags.File,
		*hostctlFlags.Sudo,
		*hostctlFlags.DryRun,
		*hostctlFlags.Name,
	)
}

// hostctlCmd represents the hostctl command
var hostctlCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "hostctl",
	Short: "modify an /etc/hosts-type file",
}

func init() { //nolint:gochecknoinits
	rootCmd.AddCommand(hostctlCmd)
	hostctlFlags.File = hostctlCmd.PersistentFlags().String("file", hostctl.DefaultHostsFile, "which hosts file to modify")
	hostctlFlags.Sudo = hostctlCmd.PersistentFlags().Bool("sudo", hostctl.DefaultSudo, "use sudo to write the hosts file")
	hostctlFlags.Name = hostctlCmd.PersistentFlags().String("name", hostctl.DefaultName, "controller name")
	hostctlFlags.DryRun = hostctlCmd.PersistentFlags().Bool("dry-run", hostctl.DefaultDryRun, "dry run")
}
