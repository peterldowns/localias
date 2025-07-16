package root

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/cmd/localias/shared"
	"github.com/peterldowns/localias/pkg/config"
)

var importCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:     "import [path]",
	Aliases: []string{"upsert"},
	Short:   "import all aliases from another config file",
	Example: shared.Example(`
# import aliases from a file named .localias.mybranch.yaml in the current directory
localias import ./.localias.mybranch.yaml

# import aliases from a file in another directory
localias import path/to/another/config/localias.yaml
	`),
	RunE: importImpl,
}

func importImpl(_ *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("invalid arguments: expected [path]")
	}
	importPath := args[0]

	// Read the config file to be imported
	importCfg, err := config.Open(importPath)
	if err != nil {
		return fmt.Errorf("failed to open import file: %w", err)
	}

	// Get the current config
	cfg := shared.Config()

	// Add/update entries from the imported config
	added, updated := cfg.Import(importCfg)

	// Save the updated config
	if err := cfg.Save(); err != nil {
		return err
	}

	for _, entry := range added {
		shared.PrintUpdate(entry, false)
	}
	for _, entry := range updated {
		shared.PrintUpdate(entry, true)
	}
	return nil
}

func init() {
	Command.AddCommand(importCmd)
}
