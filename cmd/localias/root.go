package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/pkg/config"
	"github.com/peterldowns/localias/pkg/util"
)

// These will be set at build time with ldflags, see Justfile for how they're
// defined and passed.
var (
	Version = "unknown" //nolint:gochecknoglobals
	Commit  = "unknown" //nolint:gochecknoglobals
)

var rootFlags struct { //nolint:gochecknoglobals
	Configfile *string
}

var rootCmd = &cobra.Command{ //nolint:gochecknoglobals
	Version: fmt.Sprintf("%s+commit.%s", Version, Commit),
	Use:     "localias",
	Short:   "securely proxy domains to local development servers",
	Example: util.Example(`
# Add an alias forwarding https://secure.lkl to http://127.0.0.1:9000
localias set --alias secure.lkl -p 9000
# Update an existing alias to forward to a different port
localias set --alias secure.lkl -p 9001
# Remove an alias
localias remove secure.lkl
# Show aliases
localias list
# Clear all aliases
localias clear
# Run the server, automatically applying all necessary rules to
# /etc/hosts and creating any necessary TLS certificates
localias run
# Run the server as a daemon
localias daemon start
# Check whether or not the daemon is running
localias daemon status
# Reload the config that the daemon is using
localias daemon reload
# Stop the daemon if it is running
localias daemon stop
  `),
}

func init() { //nolint:gochecknoinits
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.TraverseChildren = true
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true
	rootCmd.SetVersionTemplate("{{.Version}}\n")
	rootFlags.Configfile = rootCmd.PersistentFlags().StringP("configfile", "c", "", "path to the configuration file to edit")
}

func loadConfig() *config.Config {
	cfg, err := config.Load(rootFlags.Configfile)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}
	return cfg
}
