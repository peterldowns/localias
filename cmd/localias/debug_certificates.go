package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var certFlags struct { //nolint:gochecknoglobals
	Print *bool
}

func certImpl(_ *cobra.Command, _ []string) error {
	cfg := loadConfig()
	caddyStatePath := cfg.CaddyStatePath()
	rootCrtPath := filepath.Join(caddyStatePath, "pki/authorities/local/root.crt")
	if !*certFlags.Print {
		fmt.Println(rootCrtPath)
		return nil
	}
	content, err := os.ReadFile(rootCrtPath)
	if err != nil {
		return err
	}
	fmt.Println(string(content))
	return nil
}

var certCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "cert",
	Short: "show the certification file path",
	RunE:  certImpl,
}

func init() { //nolint:gochecknoinits
	certFlags.Print = certCmd.Flags().BoolP("print", "p", false, "print the contents of the certificate")
	debugCmd.AddCommand(certCmd)
}
