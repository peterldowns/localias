package root

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/cmd/localias/shared"
	"github.com/peterldowns/localias/pkg/config"
)

var importCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:     "import [path]...",
	Aliases: []string{"upsert"},
	Short:   "import all aliases from one or more other config files",
	Example: shared.Example(`
# import aliases from a file named .localias.mybranch.yaml in the current directory
localias import ./.localias.mybranch.yaml

# import aliases from multiple files
localias import ./.localias.mybranch.yaml ../path/to/another.localias.yaml
	`),
	RunE: importImpl,
}

func importImpl(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("invalid arguments: expected at least one [path]")
	}

	cfg := shared.Config()
	for _, importPath := range args {
		importCfg, err := config.Open(importPath)
		if err != nil {
			return fmt.Errorf("failed to open import file: %w", err)
		}

		added, updated := cfg.Import(importCfg)
		for _, entry := range added {
			shared.PrintUpdate(entry, false)
		}
		for _, entry := range updated {
			shared.PrintUpdate(entry, true)
		}
	}

	return cfg.Save()
}

func init() {
	Command.AddCommand(importCmd)
}
