package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/peterldowns/pfpro/pkg/pfpro"
)

var addFlags struct { //nolint:gochecknoglobals
	Port  *int
	Alias *string
}

func addImpl(_ *cobra.Command, _ []string) error {
	cfg, err := pfpro.Load(nil)
	if err != nil {
		return err
	}
	port := *addFlags.Port
	alias := *addFlags.Alias

	upstream := alias
	downstream := fmt.Sprintf(":%d", port)
	d := pfpro.Directive{
		Upstream:   upstream,
		Downstream: downstream,
	}
	cfg.Directives = append(cfg.Directives, d)
	if err := pfpro.WriteConfig(cfg); err != nil {
		return err
	}
	fmt.Printf(
		"%s %s -> %s\n",
		color.New(color.FgGreen).Sprint("[added]"),
		color.New(color.FgBlue).Sprint(upstream),
		color.New(color.FgWhite).Sprint(downstream),
	)
	return nil
}

var addCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "add",
	Short: "add an alias",
	RunE:  addImpl,
}

func init() { //nolint:gochecknoinits
	addFlags.Port = addCmd.Flags().IntP("port", "p", 0, "local port e.g. 9000")
	addFlags.Alias = addCmd.Flags().StringP("alias", "a", "", "domain alias e.g. example.local")
	_ = addCmd.MarkFlagRequired("port")
	_ = addCmd.MarkFlagRequired("alias")
	rootCmd.AddCommand(addCmd)
}
