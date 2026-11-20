package root

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/cmd/localias/shared"
	"github.com/peterldowns/localias/pkg/config"
)

const setPositionalArgCount = 2 // alias followed by port or host:port

var setFlags struct { //nolint:gochecknoglobals
	Port     *int
	Alias    *string
	Upstream *string
}

var setCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:     "set",
	Short:   "add or edit an alias",
	Aliases: []string{"add", "upsert", "update", "edit"},
	Example: shared.Example(`
# Add secure aliases (automatically upgrade http:// requests to https://)
## alias https://secure-explicit.test to 127.0.0.1:9001
localias set https://secure-explicit.test 9001
## alias https://secure-implicit.test to 127.0.0.1:9002
localias set secure-implicit.test 9002

# Add insecure aliases (only support http:// requests)
## alias http://not-secure.test to 127.0.0.1:9003
localias set http://not-secure.test 9003

# Add multiple aliases for the same local port
localias set door1.test 9000
localias set door2.test 9000

# Update an existing alias
localias set example.test 9001
localias set example.test 9002

# Alternative forms
localias set example.test 9001
localias set -a example.test -p 9001
localias set --alias example.test --port 9001

# Forward to a reachable host or container
localias set api.test backend:8080
localias set --alias api.test --upstream backend:8080
	`),
	RunE: setImpl,
}

func setImpl(cmd *cobra.Command, args []string) error {
	entry, err := parseSetEntry(cmd, args)
	if err != nil {
		return err
	}

	cfg, err := shared.LoadConfig()
	if err != nil {
		return err
	}
	updated := cfg.Set(entry)
	if err := cfg.Save(); err != nil {
		return err
	}
	shared.PrintUpdate(entry, updated)
	return nil
}

// parseSetEntry validates either positional arguments or explicit flags and
// normalizes the target into an Entry.
func parseSetEntry(cmd *cobra.Command, args []string) (config.Entry, error) {
	// Check whether flags were supplied, rather than whether their values differ
	// from defaults: --port=0 and --alias="" must still count as flag inputs.
	aliasFlag := cmd.Flags().Changed("alias")
	portFlag := cmd.Flags().Changed("port")
	upstreamFlag := cmd.Flags().Changed("upstream")

	// Requiring one complete input form avoids silently ignoring positional
	// arguments or conflicting flags.
	if len(args) > 0 {
		if aliasFlag || portFlag || upstreamFlag {
			return config.Entry{}, fmt.Errorf("cannot mix positional arguments with --alias, --port, or --upstream")
		}
		if len(args) != setPositionalArgCount {
			return config.Entry{}, fmt.Errorf("invalid arguments: expected [alias] [port or host:port]")
		}
		alias := args[0]
		if alias == "" {
			return config.Entry{}, fmt.Errorf("alias must not be empty")
		}
		entry, err := config.ParseTarget(args[1])
		if err != nil {
			return config.Entry{}, fmt.Errorf("invalid target: %w", err)
		}
		entry.Alias = alias
		return entry, nil
	}

	// Equal booleans mean both target flags are present or both are absent.
	if !aliasFlag || portFlag == upstreamFlag {
		return config.Entry{}, fmt.Errorf("expected --alias and exactly one of --port or --upstream")
	}
	alias := *setFlags.Alias
	if alias == "" {
		return config.Entry{}, fmt.Errorf("alias must not be empty")
	}
	if portFlag {
		entry, err := config.ParseTarget(*setFlags.Port)
		if err != nil {
			return config.Entry{}, fmt.Errorf("invalid target: %w", err)
		}
		entry.Alias = alias
		return entry, nil
	}
	// Check empty upstream explicitly because NormalizeTarget reports a missing
	// target rather than an empty value when both Port and Upstream are zero.
	if *setFlags.Upstream == "" {
		return config.Entry{}, fmt.Errorf("invalid target: target must not be empty")
	}
	entry, err := config.Entry{Alias: alias, Upstream: *setFlags.Upstream}.NormalizeTarget()
	if err != nil {
		return config.Entry{}, fmt.Errorf("invalid target: %w", err)
	}
	return entry, nil
}

func init() {
	setFlags.Alias = setCmd.Flags().StringP("alias", "a", "", "domain alias e.g. example.test")
	setFlags.Port = setCmd.Flags().IntP("port", "p", 0, "local port e.g. 9000")
	setFlags.Upstream = setCmd.Flags().String("upstream", "", "HTTP upstream address e.g. backend:8080")
	Command.AddCommand(setCmd)
}
