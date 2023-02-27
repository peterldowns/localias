package cmd

import (
	"fmt"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig"
	_ "github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/caddyauth"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/encode"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/encode/brotli"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/encode/gzip"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/encode/zstd"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/fileserver"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/headers"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/map"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/push"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/requestbody"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/reverseproxy"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/reverseproxy/fastcgi"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/reverseproxy/forwardauth"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/rewrite"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/templates"
	_ "github.com/caddyserver/caddy/v2/modules/caddypki"
	_ "github.com/caddyserver/caddy/v2/modules/caddytls"
	_ "github.com/caddyserver/caddy/v2/modules/filestorage"
	_ "github.com/caddyserver/caddy/v2/modules/logging"
	_ "github.com/caddyserver/caddy/v2/modules/metrics"
	"github.com/spf13/cobra"

	"github.com/peterldowns/pfpro/pkg/hostctl"
)

// TODO: actually daemonize like
// https://github.com/caddyserver/caddy/blob/be53e432fcac0a9b9accbc36885304639e8ca70b/cmd/commandfuncs.go#L146
// reloading should be as simple as calling caddy.Load() again
func daemonImpl(_ *cobra.Command, _ []string) error {
	hostc := controller()
	configs := Config{
		Directive{Upstream: "https://local.test", Downstream: ":9000"},
		Directive{Upstream: "https://peter.test", Downstream: ":8000"},
	}
	var added []*hostctl.Line
	for _, directive := range configs {
		x, err := httpcaddyfile.ParseAddress(directive.Upstream)
		if err != nil {
			return err
		}
		a, err := hostc.Add(true, "127.0.0.1", x.Host)
		if err != nil {
			return err
		}
		added = append(added, a...)
	}
	if added != nil {
		if err := hostc.Save(); err != nil {
			return err
		}
	}

	caddyfile := configs.Caddyfile()
	fmt.Println(caddyfile)
	cfgAdapter := caddyconfig.GetAdapter("caddyfile")
	config, warnings, err := cfgAdapter.Adapt([]byte(caddyfile), map[string]any{
		"filename": "Caddyfile",
	})
	if err != nil {
		return err
	}

	for _, w := range warnings {
		fmt.Printf("[warning]: %s\n", w.String())
	}
	err = caddy.Load(config, false)
	if err != nil {
		return err
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

// TODO: move this stuff into a separate package, not the CLI!
type (
	Config    []Directive
	Directive struct {
		Upstream   string
		Downstream string
	}
)

func (directive Directive) Caddyfile() string {
	tls := "# tls disabled"
	a, _ := httpcaddyfile.ParseAddress(directive.Upstream)
	fmt.Printf("%+v\n", a)

	if a.Scheme == "https" {
		// if strings.HasPrefix(directive.Upstream, "https://") {
		tls = strings.TrimSpace(`
	tls {
		issuer internal {
			on_demand
		}
	}
`)
	}
	return fmt.Sprintf(strings.TrimSpace(`
%s {
	reverse_proxy %s
	%s
}
	`), directive.Upstream, directive.Downstream, tls)
}

func (c Config) Caddyfile() string {
	global := strings.TrimSpace(`
{
	admin off
	persist_config off
	local_certs
	ocsp_stapling off
	storage file_system /Users/pd/.config/pfpro
	pki {
		ca local {
			name pfpro
			root_cn pfpro
			intermediate_cn pfpro
		}
	}
}
`)
	blocks := []string{global}
	for _, x := range c {
		blocks = append(blocks, x.Caddyfile())
	}
	return strings.Join(blocks, "\n") + "\n" // extra newline prevents "caddy fmt" warning in logs
}
