package main

import (
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
  `),
}

func init() { //nolint:gochecknoinits
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.TraverseChildren = true
	rootCmd.SilenceErrors = true
	// rootCmd.SilenceUsage = true
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
