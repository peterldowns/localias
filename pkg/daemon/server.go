package daemon

import (
	"fmt"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"

	"github.com/peterldowns/localias/pkg/config"
	"github.com/peterldowns/localias/pkg/hostctl"
)

func Run(hctl *hostctl.Controller, cfg *config.Config) error {
	if err := hctl.Clear(); err != nil {
		return err
	}
	for _, directive := range cfg.Directives {
		up, err := httpcaddyfile.ParseAddress(directive.Alias)
		if err != nil {
			return err
		}
		if err := hctl.Set("127.0.0.1", up.Host); err != nil {
			return err
		}
	}
	if err := hctl.Apply(); err != nil {
		return err
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
	return nil
}

func Stop() error {
	return caddy.Stop()
}
