package debug

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/peterldowns/localias/cmd/localias/shared"
	"github.com/peterldowns/localias/pkg/wsl"
)

var certFlags struct { //nolint:gochecknoglobals
	Print   *bool
	Install *bool
}

func certImpl(_ *cobra.Command, _ []string) error {
	cfg := shared.Config()
	caddyStatePath := cfg.CaddyStatePath()
	rootCrtPath := filepath.Join(caddyStatePath, "pki/authorities/local/root.crt")
	fmt.Println(rootCrtPath)
	if *certFlags.Print {
		content, err := os.ReadFile(rootCrtPath)
		if err != nil {
			return err
		}
		fmt.Println(string(content))
	}
	if *certFlags.Install {
		if err := wsl.InstallCert(rootCrtPath); err != nil {
			return err
		}
	}
	return nil
}

var certCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "cert",
	Short: "show the certification file path",
	RunE:  certImpl,
}

func init() { //nolint:gochecknoinits
	certFlags.Print = certCmd.Flags().BoolP("print", "p", false, "print the contents of the certificate")
	certFlags.Install = certCmd.Flags().BoolP("install", "i", false, "install the certificate to the windows cert store")
	Command.AddCommand(certCmd)
}
