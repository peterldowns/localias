package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var configFlags struct { //nolint:gochecknoglobals
	Print *bool
}

func configImpl(_ *cobra.Command, _ []string) error {
	cfg := loadConfig()
	if !*configFlags.Print {
		fmt.Println(cfg.Path)
		return nil
	}
	fmt.Println(cfg.Caddyfile())
	content, err := os.ReadFile(cfg.Path)
	if err != nil {
		return err
	}
	fmt.Println(string(content))
	return nil
}

var configCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "config",
	Short: "show the configuration file path",
	RunE:  configImpl,
}

func init() { //nolint:gochecknoinits
	configFlags.Print = configCmd.Flags().BoolP("print", "p", false, "print the contents of the config file")
	debugCmd.AddCommand(configCmd)
}
