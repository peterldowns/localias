package root

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/cmd/localias/shared"
	"github.com/peterldowns/localias/pkg/config"
)

var listCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:     "list",
	Aliases: []string{"l", "ls", "show"},
	Short:   "list aliases from the saved configuration",
	Long:    "List aliases from the saved configuration file. This does not query the running daemon; reload it to apply saved changes.",
	Args:    cobra.NoArgs,
	RunE:    listImpl,
}

func listImpl(cmd *cobra.Command, _ []string) error {
	format, err := cmd.Flags().GetString("format")
	if err != nil {
		return err
	}
	jsonOutput, err := cmd.Flags().GetBool("json")
	if err != nil {
		return err
	}
	if jsonOutput {
		format = "json"
	}
	if format != "text" && format != "json" {
		return fmt.Errorf("unsupported output format %q: expected text or json", format)
	}
	cfg, err := shared.LoadConfig()
	if err != nil {
		return err
	}
	w := cmd.OutOrStdout()
	switch format {
	case "json":
		return outputJSON(w, cfg.Entries)
	default:
		for _, entry := range cfg.Entries {
			if _, err := fmt.Fprintf(w,
				"%s -> %s\n",
				color.New(color.FgBlue).Sprint(entry.Alias),
				color.New(color.FgWhite).Sprintf("%d", entry.Port),
			); err != nil {
				return err
			}
		}
		return nil
	}
}

func init() {
	Command.AddCommand(listCmd)
	listCmd.Flags().Bool("json", false, "output saved aliases as JSON")
	listCmd.Flags().String("format", "text", "output format (text|json)")
	listCmd.MarkFlagsMutuallyExclusive("json", "format")
}

func aliasMap(entries []config.Entry) map[string]int {
	aliases := make(map[string]int, len(entries))
	for _, entry := range entries {
		aliases[entry.Alias] = entry.Port
	}
	return aliases
}

func outputJSON(w io.Writer, entries []config.Entry) error {
	output := struct {
		Aliases map[string]int `json:"aliases"`
	}{Aliases: aliasMap(entries)}
	if err := json.NewEncoder(w).Encode(output); err != nil {
		return fmt.Errorf("write JSON: %w", err)
	}
	return nil
}
