package debug

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/cmd/localias/shared"
	"github.com/peterldowns/localias/pkg/config"
)

var configFlags struct { //nolint:gochecknoglobals
	Print *bool
}

func caddyImpl(_ *cobra.Command, _ []string) error {
	cfg, err := shared.LoadConfig()
	if err != nil {
		return err
	}
	caddy := cfg.Caddyfile()
	fmt.Println(caddy)
	return nil
}

func configImpl(_ *cobra.Command, _ []string) error {
	path, err := config.Path(shared.Flags.Configfile)
	if err != nil {
		return fmt.Errorf("failed to find config: %w", err)
	}
	if *configFlags.Print {
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		fmt.Println(string(content))
		return nil
	}
	fmt.Println(path)
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
