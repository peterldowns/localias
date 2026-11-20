package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/peterldowns/testy/assert"

	"github.com/peterldowns/localias/pkg/hostctl"
)

func TestApplyKeepsAliasesOnLoopback(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "hosts")
	assert.NoError(t, os.WriteFile(path, []byte("127.0.0.1 localhost\n"), 0o600))
	cfg := &Config{Entries: []Entry{
		{Alias: "app.test", Port: 8000},
		{Alias: "https://api.test", Upstream: "192.168.1.10:8000"},
		{Alias: "http://frontend.test", Upstream: "frontend:5173"},
		{Alias: "ipv6.test", Upstream: "[::1]:9000"},
	}}
	assert.NoError(t, Apply(hostctl.NewFileController(path, false, "localias"), cfg))
	// Reopen the saved hosts file so this checks the persisted IPs and aliases.
	entries, err := hostctl.NewFileController(path, false, "localias").List()
	assert.NoError(t, err)
	assert.Equal(t, len(cfg.Entries), len(entries[path]))
	aliases := make(map[string]string)
	for _, line := range entries[path] {
		aliases[line.Entry.Aliases[0]] = line.Entry.IPAddress
	}
	assert.Equal(t, map[string]string{
		"app.test":      "127.0.0.1",
		"api.test":      "127.0.0.1",
		"frontend.test": "127.0.0.1",
		"ipv6.test":     "127.0.0.1",
	}, aliases)
}
