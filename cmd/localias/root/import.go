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

type importUpdate struct {
	entry   config.Entry
	updated bool
}

func importImpl(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("invalid arguments: expected at least one [path]")
	}

	cfg, err := shared.LoadConfig()
	if err != nil {
		return err
	}
	// Open and validate every import file before mutating the target config or
	// printing any output.
	importCfgs := make([]*config.Config, 0, len(args))
	for _, importPath := range args {
		importCfg, err := config.Open(importPath)
		if err != nil {
			return fmt.Errorf("failed to open import file: %w", err)
		}
		importCfgs = append(importCfgs, importCfg)
	}

	var updates []importUpdate
	for _, importCfg := range importCfgs {
		added, updated := cfg.Import(importCfg)
		for _, entry := range added {
			updates = append(updates, importUpdate{entry: entry, updated: false})
		}
		for _, entry := range updated {
			updates = append(updates, importUpdate{entry: entry, updated: true})
		}
	}

	if err := cfg.Save(); err != nil {
		return err
	}
	for _, u := range updates {
		shared.PrintUpdate(u.entry, u.updated)
	}
	return nil
}

func init() {
	Command.AddCommand(importCmd)
}
