package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "pfpro",
	Short: "securely proxy domains to local development servers",
}

func init() { //nolint:gochecknoinits
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
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
			OnError(fmt.Errorf("panic: %+v", t))
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
