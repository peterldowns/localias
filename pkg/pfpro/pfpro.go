package pfpro

import (
	"fmt"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"

	"github.com/peterldowns/pfpro/pkg/hostctl"
)

func Run(hctl *hostctl.Controller, cfg Config) error {
	var added []*hostctl.Line
	for _, directive := range cfg {
		x, err := httpcaddyfile.ParseAddress(directive.Upstream)
		if err != nil {
			return err
		}
		a, err := hctl.Add(true, "127.0.0.1", x.Host)
		if err != nil {
			return err
		}
		added = append(added, a...)
	}
	if added != nil {
		if err := hctl.Save(); err != nil {
			return err
		}
	}
	caddyfile := cfg.Caddyfile()
	cfgAdapter := caddyconfig.GetAdapter("caddyfile")
	if cfgAdapter == nil {
		panic(fmt.Errorf("nil cfgadapter"))
	}
	config, warnings, err := cfgAdapter.Adapt([]byte(caddyfile), map[string]any{
		"filename": "Caddyfile",
	})
	if err != nil {
		fmt.Printf("failed to get cfg adapter i think: %v\n", err)
		return err
	}

	for _, w := range warnings {
		fmt.Printf("[warning]: %s\n", w.String())
	}
	// TODO: actually daemonize like
	// https://github.com/caddyserver/caddy/blob/be53e432fcac0a9b9accbc36885304639e8ca70b/cmd/commandfuncs.go#L146
	// reloading should be as simple as calling caddy.Load() again
	err = caddy.Load(config, false)
	if err != nil {
		fmt.Printf("failing with %v\n", err)
		return err
	}
	select {} //nolint:revive // valid empty block
}
