package config

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/adrg/xdg"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddytls"

	// Register the same Caddy modules used by the server.
	_ "github.com/caddyserver/caddy/v2/modules/standard"
	"github.com/peterldowns/testy/assert"
)

func TestCaddyJSON(t *testing.T) { //nolint:paralleltest // isolates Caddy storage via the global XDG state directory
	previousStateHome := xdg.StateHome
	xdg.StateHome = t.TempDir()
	t.Cleanup(func() { xdg.StateHome = previousStateHome })

	cfg := Config{Entries: []Entry{
		{Alias: "bare", Port: 9000},
		{Alias: "frontend.local", Port: 9001},
		{Alias: "https://secure.test", Port: 9002},
		{Alias: "http://insecure.test", Port: 9003},
	}}
	data, _, err := cfg.CaddyJSON()
	assert.NoError(t, err)

	var caddyConfig caddy.Config
	assert.NoError(t, json.Unmarshal(data, &caddyConfig))
	assert.NotEqual(t, nil, caddyConfig.Admin)
	assert.Equal(t, true, caddyConfig.Admin.Disabled)

	var tlsConfig caddytls.TLS
	assert.NoError(t, json.Unmarshal(caddyConfig.AppsRaw["tls"], &tlsConfig))
	assert.NotEqual(t, nil, tlsConfig.Automation)
	// Internal, on-demand issuance must remain restricted to HTTPS aliases,
	// without requiring an ask endpoint or the Caddy admin listener.
	var onDemandHosts []string
	for _, policy := range tlsConfig.Automation.Policies {
		if !policy.OnDemand {
			continue
		}
		assert.NotEqual(t, 0, len(policy.SubjectsRaw))
		onDemandHosts = append(onDemandHosts, policy.SubjectsRaw...)
		assert.Equal(t, 1, len(policy.IssuersRaw))
		var issuer struct {
			Module string `json:"module"`
			CA     string `json:"ca"`
		}
		assert.NoError(t, json.Unmarshal(policy.IssuersRaw[0], &issuer))
		assert.Equal(t, "internal", issuer.Module)
		assert.Equal(t, "local", issuer.CA)
	}
	slices.Sort(onDemandHosts)
	assert.Equal(t, []string{"bare", "frontend.local", "secure.test"}, onDemandHosts)
	assert.Equal(t, (*caddytls.OnDemandConfig)(nil), tlsConfig.Automation.OnDemand)
}
