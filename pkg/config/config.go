package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"go.yaml.in/yaml/v2"
)

// Config files are readable by other users, but writable only by their owner.
const configFileMode = 0o644

type Config struct {
	Path    string
	Entries []Entry
}

// Entry maps an alias to an HTTP target. Exactly one of Port or Upstream is set.
type Entry struct {
	Alias    string
	Port     int
	Upstream string
}

// Set will add or update the existing list of entries.  If there is
// already an entry with the same alias, its target will be
// updated. If there is not already an entry with the same alias,
// the new entry will be added to the list.
//
// Returns `true` if an existing entry was updated, `false` if the entry
// was added.
func (c *Config) Set(d Entry) bool {
	for i, existing := range c.Entries {
		if existing.Alias == d.Alias {
			c.Entries[i] = d
			return true
		}
	}
	c.Entries = append(c.Entries, d)
	return false
}

// Import merges the entries from another config into the current one.
// It returns the entries that were added and updated.
func (c *Config) Import(other *Config) (added []Entry, updated []Entry) {
	for _, entry := range other.Entries {
		if c.Set(entry) {
			updated = append(updated, entry)
		} else {
			added = append(added, entry)
		}
	}
	return added, updated
}

// Remove removes all entries from the config that match
// any of the specified aliases.
func (c *Config) Remove(aliases ...string) []Entry {
	var removed []Entry
	previous := c.Entries
	c.Clear()
	for _, d := range previous {
		shouldRemove := false
		for _, alias := range aliases {
			if d.Alias == alias {
				shouldRemove = true
				break
			}
		}
		if shouldRemove {
			removed = append(removed, d)
		} else {
			c.Set(d)
		}
	}
	return removed
}

// Clear removes all entries from the config.
func (c *Config) Clear() []Entry {
	removed := c.Entries
	c.Entries = []Entry{}
	return removed
}

// Save writes the config to disk.
func (c *Config) Save() error {
	// Build an ordered YAML mapping with integers for local ports and strings
	// for addresses. Validate every entry before writing so a programmatically
	// constructed invalid target cannot overwrite the saved configuration.
	entries := yaml.MapSlice{}
	for _, d := range c.Entries {
		normalized, err := d.NormalizeTarget()
		if err != nil {
			return d.invalidTargetError(err)
		}
		d = normalized
		var value any = d.Port
		if d.Upstream != "" {
			value = d.Upstream
		}
		entries = append(entries, yaml.MapItem{
			Key:   d.Alias,
			Value: value,
		})
	}
	bytes := []byte(strings.TrimSpace(`
# localias config file syntax
#
#   alias: port or "host:port"
#
# for example,
#
#   bareTLD: 9003 # serves over https and http
#   implicitly_secure.test: 9002 # serves over https and http
#   https://explicit_secure.test: 9000 # serves over https and http
#   http://explicit_insecure.test: 9001 # serves over http only
#   api.test: "192.168.1.10:8000" # HTTP upstream at an explicit address
#   frontend.test: "frontend:5173" # HTTP upstream resolved by Caddy
#
# Ports use localhost on the machine or container running localias.
#
	`) + "\n")
	if len(entries) != 0 {
		entryBytes, err := yaml.Marshal(entries)
		if err != nil {
			return err
		}
		bytes = append(bytes, entryBytes...)
	}
	return os.WriteFile(c.Path, bytes, configFileMode)
}

func (Config) CaddyStatePath() string {
	path, err := xdg.StateFile("localias/caddy")
	if err != nil {
		panic(err)
	}
	path, err = filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return path
}

func (c Config) Caddyfile() string {
	return c.caddyfile(false)
}

func (c Config) caddyfile(normalized bool) string {
	path := c.CaddyStatePath()
	allowedMap := ""
	for _, entry := range c.Entries {
		allowedMap += fmt.Sprintf("		%s 1\n", entry.Host())
	}
	global := fmt.Sprintf(strings.TrimSpace(`
{
	admin off
	persist_config off
	local_certs
	ocsp_stapling off
	storage file_system "%s"

	pki {
		ca local {
			name "Localias"
			root_cn "Localias Root"
			intermediate_cn "Localias Intermediate"
		}
	}

}
`), path)
	blocks := []string{global}
	for _, x := range c.Entries {
		if normalized {
			blocks = append(blocks, x.caddyfileWithUpstream(x.rawUpstreamAddress()))
		} else {
			blocks = append(blocks, x.Caddyfile())
		}
	}
	// extra newline prevents "caddy fmt" warning in logs
	return strings.Join(blocks, "\n") + "\n"
}

func (c Config) CaddyJSON() ([]byte, []caddyconfig.Warning, error) {
	// Callers can construct entries without going through the YAML or CLI parser.
	// Revalidate before interpolating targets into Caddy's configuration syntax.
	// Copy the slice before normalizing to avoid changing the caller's entries.
	entries := make([]Entry, len(c.Entries))
	for i, entry := range c.Entries {
		normalized, err := entry.NormalizeTarget()
		if err != nil {
			return nil, nil, entry.invalidTargetError(err)
		}
		entries[i] = normalized
	}
	normalizedConfig := c
	normalizedConfig.Entries = entries
	caddyfile := normalizedConfig.caddyfile(true)
	cfgAdapter := caddyconfig.GetAdapter("caddyfile")
	if cfgAdapter == nil {
		return nil, nil, fmt.Errorf("failed to load caddyfile adapater")
	}
	cfgJSON, warnings, err := cfgAdapter.Adapt([]byte(caddyfile), map[string]any{
		"filename": "Caddyfile",
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse configuration: %w", err)
	}
	return cfgJSON, warnings, nil
}

func (entry Entry) String() string {
	return fmt.Sprintf("%s: %s", entry.Alias, entry.Target())
}

func (entry Entry) Host() string {
	a, _ := httpcaddyfile.ParseAddress(entry.Alias)
	return a.Host
}

func (entry Entry) Caddyfile() string {
	return entry.caddyfileWithUpstream(entry.UpstreamAddress())
}

func (entry Entry) caddyfileWithUpstream(upstream string) string {
	tls := "# tls disabled"
	a, _ := httpcaddyfile.ParseAddress(entry.Alias)
	// If no scheme is given, default to https.
	if a.Scheme == "" {
		a.Scheme = "https"
	}
	if a.Scheme == "https" {
		tls = strings.TrimSpace(`
	tls {
		on_demand
		issuer internal {
			ca local
		}
	}
`)
	}
	return fmt.Sprintf(strings.TrimSpace(`
%s {
	reverse_proxy %s
	%s
}
	`), entry.Alias, upstream, tls)
}
