package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig"             // config
	_ "github.com/caddyserver/caddy/v2/caddyconfig/caddyfile" // caddy?
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp/reverseproxy"
	"github.com/caddyserver/caddy/v2/modules/caddytls"
	_ "github.com/caddyserver/caddy/v2/modules/standard" // ahhhhhhhh?

	"github.com/spf13/cobra"
)

func daemonImpl(_ *cobra.Command, _ []string) error {
	// potentially: build config as caddyfile, much simler?
	// https://caddyserver.com/docs/caddyfile/directives/tls#internal

	// proceed to build the handler and server
	// based on
	// https://github.com/caddyserver/caddy/blob/5ded580444e9258cb35a9c94192d3c1d63e7b74f/modules/caddyhttp/reverseproxy/command.go#L122
	ht := reverseproxy.HTTPTransport{}
	toAddresses := []string{":9000"}
	upstreamPool := reverseproxy.UpstreamPool{}
	for _, toAddr := range toAddresses {
		upstreamPool = append(upstreamPool, &reverseproxy.Upstream{
			Dial: toAddr,
		})
	}
	handler := reverseproxy.Handler{
		TransportRaw: caddyconfig.JSONModuleObject(ht, "protocol", "http", nil),
		Upstreams:    upstreamPool,
	}
	route := caddyhttp.Route{
		HandlersRaw: []json.RawMessage{
			caddyconfig.JSONModuleObject(handler, "handler", "reverse_proxy", nil),
		},
	}
	fromAddr, err := httpcaddyfile.ParseAddress("https://local.test")
	if err != nil {
		return err
	}
	fromAddr.Host = "local.test"
	fromAddr.Scheme = "https"
	fromAddr.Port = "443"
	if fromAddr.Host != "" {
		route.MatcherSetsRaw = []caddy.ModuleMap{
			{
				"host": caddyconfig.JSON(caddyhttp.MatchHost{fromAddr.Host}, nil),
			},
		}
	}
	server := &caddyhttp.Server{
		Routes: caddyhttp.RouteList{route},
		Listen: []string{":" + fromAddr.Port},
	}
	// server.AutoHTTPS = &caddyhttp.AutoHTTPSConfig{DisableRedir: true}
	httpApp := caddyhttp.App{
		Servers: map[string]*caddyhttp.Server{"proxy": server},
	}
	appsRaw := caddy.ModuleMap{
		"http": caddyconfig.JSON(httpApp, nil),
	}
	tlsApp := caddytls.TLS{
		Automation: &caddytls.AutomationConfig{
			Policies: []*caddytls.AutomationPolicy{{
				Subjects:   []string{fromAddr.Host},
				IssuersRaw: []json.RawMessage{json.RawMessage(`{"module":"internal"}`)},
			}},
		},
	}
	appsRaw["tls"] = caddyconfig.JSON(tlsApp, nil)
	var false bool
	cfg := &caddy.Config{
		Admin: &caddy.AdminConfig{
			Disabled: true,
			Config: &caddy.ConfigSettings{
				Persist: &false,
			},
		},
		AppsRaw: appsRaw,
	}
	cfg.Logging = &caddy.Logging{
		Logs: map[string]*caddy.CustomLog{
			"default": {Level: "DEBUG"},
		},
	}
	err = caddy.Run(cfg)
	if err != nil {
		return err
	}

	for _, to := range toAddresses {
		fmt.Printf("Caddy proxying %s -> %s\n", fromAddr.String(), to)
	}
	if len(toAddresses) > 1 {
		fmt.Println("Load balancing policy: random")
	}

	select {}
}

var daemonCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "daemon",
	Short: "run the caddy daemon server",
	RunE:  daemonImpl,
}

func init() { //nolint:gochecknoinits
	rootCmd.AddCommand(daemonCmd)
}
