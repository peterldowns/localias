package main

import (
	"fmt"
	"strconv"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/pkg/config"
	"github.com/peterldowns/localias/pkg/util"
)

var setFlags struct { //nolint:gochecknoglobals
	Port  *int
	Alias *string
}

func setImpl(_ *cobra.Command, args []string) error {
	cfg, err := config.Load(nil)
	if err != nil {
		return err
	}
	alias := *setFlags.Alias
	port := *setFlags.Port

	if port == 0 && alias == "" {
		if len(args) != 2 {
			return fmt.Errorf("invalid arguments: expected [alias] [port]")
		}
		alias = args[0]
		x, err := strconv.ParseInt(args[1], 0, 0)
		if err != nil {
			return fmt.Errorf("valid to parse port: %w", err)
		}
		port = int(x)
	}

	updated := cfg.Upsert(config.Entry{
		Alias: alias,
		Port:  port,
	})
	if err := cfg.Save(); err != nil {
		return err
	}

	action := "[added]"
	if updated {
		action = "[updated]"
	}
	fmt.Printf(
		"%s %s -> %s\n",
		color.New(color.FgGreen).Sprint(action),
		color.New(color.FgBlue).Sprintf(alias),
		color.New(color.FgWhite).Sprintf("%d", port),
	)
	return nil
}

var setCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:     "set",
	Short:   "add or edit an alias",
	Aliases: []string{"add", "upsert", "update", "edit"},
	Example: util.Example(`
# Add secure aliases (automatically upgrade http:// requests to https://)
## alias https://secure-explicit.lkl to 127.0.0.1:9001
localias set --alias https://secure-explicit.lkl --port 9001
## alias https://secure-implicit.lkl to 127.0.0.1:9002
localias set --alias secure-implicit.lkl --port 9002

# Add insecure aliases (only support http:// requests)
## alias http://not-secure.lkl to 127.0.0.1:9003
localias set --alias http://not-secure.lkl --port 9003

# Add multiple aliases for the same local port
localias set --alias door1.lkl --port 9000
localias set --alias door2.lkl --port 9000

# Update an existing alias
localias set --alias example.lkl --port 9001
localias set --alias example.lkl --port 9002

# Alternative forms
localias set example.lkl 9001
localias set -a example.lkl -p 9001
localias set --alias example.lkl --port 9001


	`),
	RunE: setImpl,
}

func init() { //nolint:gochecknoinits
	setFlags.Alias = setCmd.Flags().StringP("alias", "a", "", "domain alias e.g. example.lkl")
	setFlags.Port = setCmd.Flags().IntP("port", "p", 0, "local port e.g. 9000")
	setCmd.MarkFlagsRequiredTogether("alias", "port")
	rootCmd.AddCommand(setCmd)
}
