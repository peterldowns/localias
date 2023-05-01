package server

import (
	"github.com/caddyserver/caddy/v2"

	_ "github.com/peterldowns/localias/pkg/caddymodules" // necessary caddy configuration
	"github.com/peterldowns/localias/pkg/config"
	"github.com/peterldowns/localias/pkg/hostctl"
)

// Start will start the caddy server (if it hasn't been started already) and apply
// the latest configuration.
func Start(hctl hostctl.Controller, cfg *config.Config) error {
	err := config.Apply(hctl, cfg)
	if err != nil {
		return err
	}
	cfgJSON, _, err := cfg.CaddyJSON()
	if err != nil {
		return err
	}
	return caddy.Load(cfgJSON, false)
}

// Stop the caddy server (if it is running)
func Stop() error {
	return caddy.Stop()
}
