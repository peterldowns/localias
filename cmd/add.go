package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/peterldowns/pfpro/pkg/pfpro"
)

func addImpl(_ *cobra.Command, args []string) error {
	cfg, err := pfpro.Load(nil)
	if err != nil {
		return err
	}

	if len(args) < 2 {
		return fmt.Errorf("must pass at least 2 args")
	}
	port := args[0]
	aliases := args[1:]
	added := []pfpro.Directive{}
	for _, alias := range aliases {
		d := pfpro.Directive{
			Upstream:   alias,
			Downstream: port,
		}
		cfg.Directives = append(cfg.Directives, d)
		added = append(added, d)
	}
	if err := pfpro.WriteConfig(cfg); err != nil {
		return err
	}
	for _, d := range added {
		fmt.Printf(
			"%s %s -> %s\n",
			color.New(color.FgGreen).Sprint("[added]"),
			color.New(color.FgBlue).Sprint(d.Upstream),
			color.New(color.FgWhite).Sprint(d.Downstream),
		)
	}
	return nil
}

var addCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "add port domain [...more domains]",
	Args:  cobra.MinimumNArgs(2),
	Short: "add an alias",
	RunE:  addImpl,
}

func init() { //nolint:gochecknoinits
	rootCmd.AddCommand(addCmd)
}
