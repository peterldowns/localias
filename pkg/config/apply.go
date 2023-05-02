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
		if err := hctl.Set("127.0.0.1", up.Host); err != nil {
			return err
		}
		fmt.Println(entry.String())
	}
	fmt.Println("entries!")
	entries, err := hctl.List()
	if err != nil {
		panic(err)
	}
	for _, e := range entries {
		fmt.Println(e.Raw)
	}
	return hctl.Apply()
}
