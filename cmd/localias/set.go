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

	cfg := loadConfig()
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
## alias https://secure-explicit.test to 127.0.0.1:9001
localias set --alias https://secure-explicit.test --port 9001
## alias https://secure-implicit.test to 127.0.0.1:9002
localias set --alias secure-implicit.test --port 9002

# Add insecure aliases (only support http:// requests)
## alias http://not-secure.test to 127.0.0.1:9003
localias set --alias http://not-secure.test --port 9003

# Add multiple aliases for the same local port
localias set --alias door1.test --port 9000
localias set --alias door2.test --port 9000

# Update an existing alias
localias set --alias example.test --port 9001
localias set --alias example.test --port 9002

# Alternative forms
localias set example.test 9001
localias set -a example.test -p 9001
localias set --alias example.test --port 9001


	`),
	RunE: setImpl,
}

func init() { //nolint:gochecknoinits
	setFlags.Alias = setCmd.Flags().StringP("alias", "a", "", "domain alias e.g. example.test")
	setFlags.Port = setCmd.Flags().IntP("port", "p", 0, "local port e.g. 9000")
	setCmd.MarkFlagsRequiredTogether("alias", "port")
	rootCmd.AddCommand(setCmd)
}
