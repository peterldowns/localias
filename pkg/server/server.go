package server

import (
	"github.com/caddyserver/caddy/v2"

	_ "github.com/peterldowns/localias/pkg/caddymodules" // necessary caddy configuration
	"github.com/peterldowns/localias/pkg/config"
)

// Start will start the caddy server (if it hasn't been started already).
func Start(cfg *config.Config) error {
	cfgJSON, _, err := cfg.CaddyJSON()
	if err != nil {
		return err
	}
	return caddy.Load(cfgJSON, false)
}

// Stop will stop the caddy server (if it is running).
func Stop() error {
	return caddy.Stop()
}
