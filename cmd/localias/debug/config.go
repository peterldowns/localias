package debug

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/cmd/localias/shared"
)

var configFlags struct { //nolint:gochecknoglobals
	Print *bool
}

func caddyImpl(_ *cobra.Command, _ []string) error {
	cfg := shared.Config()
	caddy := cfg.Caddyfile()
	fmt.Println(caddy)
	return nil
}

func configImpl(_ *cobra.Command, _ []string) error {
	cfg := shared.Config()
	if *configFlags.Print {
		content, err := os.ReadFile(cfg.Path)
		if err != nil {
			return err
		}
		fmt.Println(string(content))
		return nil
	}
	fmt.Println(cfg.Path)
	return nil
}

var configCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "config",
	Short: "show the configuration file path",
	RunE:  configImpl,
}

var caddyCmd = &cobra.Command{
	Use:   "caddyfile",
	Short: "show the Caddy configuration file used by localias",
	RunE:  caddyImpl,
}

func init() {
	configFlags.Print = configCmd.Flags().BoolP("print", "p", false, "print the contents of the config file")
	Command.AddCommand(configCmd)
	Command.AddCommand(caddyCmd)
}
