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

type Config struct {
	Upstream      string
	Wrapper       string
	WrapperHost   string
	WrapperPort   string
	WrapperScheme string
}

func daemonImpl(_ *cobra.Command, _ []string) error {
	// potentially: build config as caddyfile, much simler?
	// https://caddyserver.com/docs/caddyfile/directives/tls#internal

	// proceed to build the handler and server
	// based on
	// https://github.com/caddyserver/caddy/blob/5ded580444e9258cb35a9c94192d3c1d63e7b74f/modules/caddyhttp/reverseproxy/command.go#L122

	cfgs := []Config{
		{WrapperPort: "443", WrapperScheme: "https", WrapperHost: "local.test", Wrapper: "https://local.test", Upstream: ":9000"},
		{WrapperPort: "80", WrapperScheme: "http", WrapperHost: "insecure.test", Wrapper: "http://inseure.test", Upstream: ":9000"},
	}
	// var servers []*caddyhttp.Server
	var routes []caddyhttp.Route
	var listeners []string
	for _, cfg := range cfgs {
		ht := reverseproxy.HTTPTransport{}
		upstreamPool := reverseproxy.UpstreamPool{
			&reverseproxy.Upstream{
				Dial: cfg.Upstream,
			},
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
		fromAddr, err := httpcaddyfile.ParseAddress(cfg.Wrapper)
		if err != nil {
			return err
		}
		fromAddr.Host = cfg.WrapperHost
		fromAddr.Scheme = cfg.WrapperScheme
		fromAddr.Port = cfg.WrapperPort
		if fromAddr.Host != "" {
			route.MatcherSetsRaw = []caddy.ModuleMap{
				{
					"host": caddyconfig.JSON(caddyhttp.MatchHost{fromAddr.Host}, nil),
				},
			}
		}
		routes = append(routes, route)
		listeners = append(listeners, ":"+fromAddr.Port)
		// server := &caddyhttp.Server{
		// 	Routes: caddyhttp.RouteList{route},
		// 	Listen: []string{":" + fromAddr.Port},
		// }
		// servers = append(servers, server)
	}
	// smap := map[string]*caddyhttp.Server{}
	// for i, server := range servers {
	// 	smap[fmt.Sprintf("proxy-%d", i)] = server
	// }
	// server.AutoHTTPS = &caddyhttp.AutoHTTPSConfig{DisableRedir: true}
	server := &caddyhttp.Server{
		Routes: routes,
		Listen: listeners,
	}
	fmt.Println("listeners", listeners)
	httpApp := caddyhttp.App{
		Servers: map[string]*caddyhttp.Server{
			"proxy": server,
		},
		// Servers: map[string]*caddyhttp.Server{"proxy": server},
	}
	appsRaw := caddy.ModuleMap{
		"http": caddyconfig.JSON(httpApp, nil),
	}
	var subjects []string
	// var policies []*caddytls.AutomationPolicy
	for _, cfg := range cfgs {
		subjects = append(subjects, cfg.WrapperHost)
		// policies = append(policies, &caddytls.AutomationPolicy{
		// 	Subjects:   []string{cfg.WrapperHost},
		// 	IssuersRaw: []json.RawMessage{json.RawMessage(`{"module":"internal"}`)},
		// })
	}
	tlsApp := caddytls.TLS{
		Automation: &caddytls.AutomationConfig{
			Policies: []*caddytls.AutomationPolicy{{
				// Subjects:   []string{fromAddr.Host},
				Subjects:   subjects,
				IssuersRaw: []json.RawMessage{json.RawMessage(`{"module":"internal"}`)},
			}},
		},
	}
	// tlsApp := caddytls.TLS{
	// 	Automation: &caddytls.AutomationConfig{
	// 		Policies: policies,
	// 	},
	// }
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
	err := caddy.Run(cfg)
	if err != nil {
		return err
	}

	for _, cfg := range cfgs {
		fmt.Printf("Caddy proxying %s -> %s\n", cfg.Wrapper, cfg.Upstream)
	}

	select {} //nolint:revive // valid empty block
}

var daemonCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "daemon",
	Short: "run the caddy daemon server",
	RunE:  daemonImpl,
}

func init() { //nolint:gochecknoinits
	rootCmd.AddCommand(daemonCmd)
}
