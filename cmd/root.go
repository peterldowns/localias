package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "localias",
	Short: "securely proxy domains to local development servers",
	Example: trimLeading(`
# Add an alias forwarding https://secure.local to http://127.0.0.1:9000
localias add --alias secure.local -p 9000
# Remove an alias
localias remove secure.local
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
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.TraverseChildren = true
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true
}

func Execute() {
	defer func() {
		switch t := recover().(type) {
		case error:
			OnError(fmt.Errorf("panic: %w", t))
		case string:
			OnError(fmt.Errorf("panic: %s", t))
		default:
			if t != nil {
				OnError(fmt.Errorf("panic: %+v", t))
			}
		}
	}()
	if err := rootCmd.Execute(); err != nil {
		OnError(err)
	}
}

func OnError(err error) {
	msg := color.New(color.FgRed, color.Italic).Sprintf("error: %s\n", err)
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

// trimLeading removes any surrounding space from a string, then removes any
// leading whitespace from each line in the string.
func trimLeading(s string) string {
	in := strings.Split(strings.TrimSpace(s), "\n")
	var out []string

	for _, x := range in {
		x = strings.TrimSpace(x)
		if len(x) > 0 && x[0] == '#' {
			x = color.New(color.Faint).Sprint(x)
		}
		out = append(out, "  "+x)
	}
	return strings.Join(out, "\n")
}
