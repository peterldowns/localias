package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "pfpro",
	Short: "securely proxy domains to local development servers",
	Example: trimLeading(`
# Add a forwarding rule: https://secure.local to http://127.0.0.1:9000
pfpro add --alias secure.local -p 9000
# Remove a forwarding rule
pfpro remove secure.local
# Show forwarding rules
pfpro list
# Clear all forwarding rules
pfpro clear

# Run the server, automatically applying all necessary rules to
# /etc/hosts and creating any necessary TLS certificates
pfpro run
	`),
}

func init() { //nolint:gochecknoinits
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.TraverseChildren = true
	rootCmd.SilenceErrors = true
	// rootCmd.SilenceUsage = true
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
