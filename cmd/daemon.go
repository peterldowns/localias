package cmd

import (
	"fmt"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig"             // config
	_ "github.com/caddyserver/caddy/v2/caddyconfig/caddyfile" // caddy?
	_ "github.com/caddyserver/caddy/v2/modules/standard"      // ahhhhhhhh?

	"github.com/spf13/cobra"
)

func daemonImpl(_ *cobra.Command, _ []string) error {
	// run the initial config
	cfgAdapter := caddyconfig.GetAdapter("caddyfile")
	if cfgAdapter == nil {
		panic("whtf")
	}
	configContents := strings.TrimSpace(`
local.test

reverse_proxy :9000
`)
	config, warnings, err := cfgAdapter.Adapt([]byte(configContents), map[string]any{
		"filename": "Caddyfile",
	})
	if err != nil {
		return err
	}
	for _, w := range warnings {
		fmt.Printf("warning; %+v\n", w)
	}
	// config := &caddy.Config{}
	// tlsApp := caddytls.TLS{
	// 	Automation: &caddytls.AutomationConfig{
	// 		Policies: []*caddytls.AutomationPolicy{{
	// 			Subjects:   []string{fromAddr.Host},
	// 			IssuersRaw: []json.RawMessage{json.RawMessage(`{"module":"internal"}`)},
	// 		}},
	// 	},
	// }
	// appsRaw["tls"] = caddyconfig.JSON(tlsApp, nil)

	if err := caddy.Load(config, true); err != nil {
		return err
	}
	// time.Sleep(30 * time.Second)
	return nil
}

var daemonCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "daemon",
	Short: "run the caddy daemon server",
	RunE:  daemonImpl,
}

func init() { //nolint:gochecknoinits
	rootCmd.AddCommand(daemonCmd)
}
