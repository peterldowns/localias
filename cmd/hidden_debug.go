package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/pkg/config"
)

func debugImpl(_ *cobra.Command, _ []string) error {
	path, err := config.DefaultPath()
	if err != nil {
		return err
	}
	fmt.Println(path)
	if err := listImpl(nil, nil); err != nil {
		return err
	}
	if err := hostctlListImpl(nil, nil); err != nil {
		return err
	}
	cfg, err := config.Load(nil)
	if err != nil {
		return err
	}
	fmt.Println(cfg.Caddyfile())
	return nil
}

var debugCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:     "debug",
	Aliases: []string{"l"},
	Short:   "debug the configuration",
	RunE:    debugImpl,
	Hidden:  true,
}

func init() { //nolint:gochecknoinits
	rootCmd.AddCommand(debugCmd)
}
