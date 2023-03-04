package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/pkg/config"
)

var addFlags struct { //nolint:gochecknoglobals
	Port  *int
	Alias *string
}

func addImpl(_ *cobra.Command, _ []string) error {
	cfg, err := config.Load(nil)
	if err != nil {
		return err
	}
	port := *addFlags.Port
	alias := *addFlags.Alias

	upstream := alias
	d := config.Directive{
		Alias: upstream,
		Port:  port,
	}
	cfg.Directives = append(cfg.Directives, d)
	if err := cfg.Save(); err != nil {
		return err
	}
	fmt.Printf(
		"%s %s -> %s\n",
		color.New(color.FgGreen).Sprint("[added]"),
		color.New(color.FgBlue).Sprint(upstream),
		color.New(color.FgWhite).Sprintf(":%d", port),
	)
	return nil
}

var addCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "add",
	Short: "add an alias",
	Example: trimLeading(`
# Add secure aliases (automatically upgrade http:// requests to https://)
## alias https://secure-explicit.local to 127.0.0.1:9001
localias add --alias https://secure-explicit.local --port 9001
## alias https://secure-implicit.local to 127.0.0.1:9002
localias add --alias secure-implicit.local --port 9002

# Add insecure aliases (only support http:// requests)
## alias http://not-secure.local to 127.0.0.1:9003
localias add --alias http://not-secure.local --port 9003

# Add multiple aliases for the same local port
localias add --alias door1.local --port 9000
localias add --alias door2.local --port 9000

# Overwrite an existing alias
localias add --alias example.local --port 9001
localias add --alias example.local --port 9002
	`),
	RunE: addImpl,
}

func init() { //nolint:gochecknoinits
	addFlags.Port = addCmd.Flags().IntP("port", "p", 0, "local port e.g. 9000")
	addFlags.Alias = addCmd.Flags().StringP("alias", "a", "", "domain alias e.g. example.local")
	_ = addCmd.MarkFlagRequired("port")
	_ = addCmd.MarkFlagRequired("alias")
	rootCmd.AddCommand(addCmd)
}
