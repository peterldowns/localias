package daemon

import (
	"fmt"

	"github.com/caddyserver/caddy/v2/caddyconfig"
	caddycmd "github.com/caddyserver/caddy/v2/cmd"

	_ "github.com/peterldowns/localias/pkg/caddymodules" // necessary caddy configuration
	"github.com/peterldowns/localias/pkg/config"
)

// Get the URI address for the locally-running caddy server based on the current
// configuration file.
//
// Based on
// https://github.com/caddyserver/caddy/blob/be53e432fcac0a9b9accbc36885304639e8ca70b/cmd/commandfuncs.go#L146
// but parses our generated config file instead of using Caddy's helpers for
// loading configs from disk.
func determineAPIAddress(cfg *config.Config) (string, error) {
	caddyfile := cfg.Caddyfile()
	cfgAdapter := caddyconfig.GetAdapter("caddyfile")
	if cfgAdapter == nil {
		return "", fmt.Errorf("failed to load caddyfile adapater")
	}
	cfgJSON, _, err := cfgAdapter.Adapt([]byte(caddyfile), map[string]any{
		"filename": "Caddyfile",
	})
	if err != nil {
		return "", fmt.Errorf("failed to parse configuration: %w", err)
	}
	address, err := caddycmd.DetermineAdminAPIAddress("", cfgJSON, caddyfile, "caddyfile")
	if err != nil {
		return "", fmt.Errorf("could not determine api address: %w", err)
	}
	return address, nil
}
