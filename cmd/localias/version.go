package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func versionImpl(_ *cobra.Command, _ []string) error {
	fmt.Printf("%s\n", rootCmd.Version)
	return nil
}

var versionCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "version",
	Short: "show the version of this binary",
	RunE:  versionImpl,
}

func init() { //nolint:gochecknoinits
	rootCmd.AddCommand(versionCmd)
}
