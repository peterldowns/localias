package config

import (
	"path/filepath"
	"testing"

	"github.com/caddyserver/caddy/v2/caddyconfig"
)

// FuzzTarget exercises the complete accepted-target path without starting Caddy
// or resolving upstream names. Adapt only the site block to avoid storage setup.
func FuzzTarget(f *testing.F) {
	for _, seed := range []string{
		"", "3000", "0003000", "frontend:5173", "localhost:80", "127.0.0.1:65535",
		"[::1]:8080", "bad:0", "http://frontend:5173", "a {\n respond hacked\n}:80",
		"[fe80::1%eth0]:80", "frontend.:1234", "a_b:3000", "\x00:80",
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, input string) {
		entry, err := ParseTarget(input)
		if err != nil {
			return
		}
		if err := entry.ValidateTarget(); err != nil {
			t.Fatalf("accepted target %q failed validation: %v", input, err)
		}
		normalized, err := ParseTarget(entry.Target())
		if err != nil || normalized != entry {
			t.Fatalf("reparse target %q: got %+v, %v; want %+v", entry.Target(), normalized, err, entry)
		}
		entry.Alias = "http://proxy.test:18080"
		cfg := Config{Path: filepath.Join(t.TempDir(), "config.yaml"), Entries: []Entry{entry}}
		if err := cfg.Save(); err != nil {
			t.Fatalf("save target %q: %v", input, err)
		}
		reopened, err := Open(cfg.Path)
		if err != nil {
			t.Fatalf("reopen target %q: %v", input, err)
		}
		if len(reopened.Entries) != 1 || reopened.Entries[0] != entry {
			t.Fatalf("reopened target %q: got %+v; want %+v", input, reopened.Entries, entry)
		}
		adapter := caddyconfig.GetAdapter("caddyfile")
		if adapter == nil {
			t.Fatal("Caddyfile adapter not registered")
		}
		if _, _, err := adapter.Adapt([]byte(entry.Caddyfile()+"\n"), nil); err != nil {
			t.Fatalf("adapt target %q: %v", input, err)
		}
	})
}
