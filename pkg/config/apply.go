package config

import (
	"fmt"

	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"

	"github.com/peterldowns/localias/pkg/hostctl"
)

func Apply(hctl hostctl.Controller, cfg *Config) error {
	if err := hctl.Clear(); err != nil {
		return err
	}
	for _, entry := range cfg.Entries {
		up, err := httpcaddyfile.ParseAddress(entry.Alias)
		if err != nil {
			return err
		}
		if err := hctl.SetLocal(up.Host); err != nil {
			return err
		}
		fmt.Println(entry.String())
	}
	return hctl.Apply()
}
