package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/go-yaml/yaml"
)

type Config struct {
	Path    string
	Entries []Entry
}

type Entry struct {
	Alias string
	Port  int
}

// Set will add or update the existing list of entries.  If there is
// already a entry with the same upstream/alias, its downstream/port will be
// updated. If there is not already a entry with the same upstream/alias,
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
	entries := yaml.MapSlice{}
	for _, d := range c.Entries {
		entries = append(entries, yaml.MapItem{
			Key:   d.Alias,
			Value: d.Port,
		})
	}
	bytes := []byte(strings.TrimSpace(`
# localias config file syntax
#
# 	alias: port
#
# for example,
#
#   bareTLD: 9003 # serves over https and http
#   implicitly_secure.test: 9002 # serves over https and http
#   https://explicit_secure.test: 9000 # serves over https and http
#   http://explicit_insecure.test: 9001 # serves over http only
#
	`) + "\n")
	if len(entries) != 0 {
		entryBytes, err := yaml.Marshal(entries)
		if err != nil {
			return err
		}
		bytes = append(bytes, entryBytes...)
	}
	return os.WriteFile(c.Path, bytes, 0o644)
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
	path := c.CaddyStatePath()
	allowedMap := ""
	for _, entry := range c.Entries {
		allowedMap += fmt.Sprintf("		%s 1\n", entry.Host())
	}
	global := fmt.Sprintf(strings.TrimSpace(`
{
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

	# Allow the internal CA to re-issue certificates for all of the sites
	# in this file at the same time, every second. We have to put a configuration
	# here or the logs will show a scary security warning.
	# https://caddyserver.com/docs/automatic-https#on-demand-tls
	on_demand_tls {
		interval 1s
		burst %d
		ask http://127.0.0.1:2019
	}
}
:2019 {
	bind 127.0.0.1 ::1

	map {query.domain} {allowed} {
		%s
		default 0
	}

	@allowed %s{allowed} == "1"%s
	respond @allowed 200
	respond 400
}
`), path, len(c.Entries)+1, allowedMap, "`", "`")
	blocks := []string{global}
	for _, x := range c.Entries {
		blocks = append(blocks, x.Caddyfile())
	}
	// extra newline prevents "caddy fmt" warning in logs
	return strings.Join(blocks, "\n") + "\n"
}

func (c Config) CaddyJSON() ([]byte, []caddyconfig.Warning, error) {
	caddyfile := c.Caddyfile()
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
	return fmt.Sprintf("%s: %d", entry.Alias, entry.Port)
}

func (entry Entry) Host() string {
	a, _ := httpcaddyfile.ParseAddress(entry.Alias)
	return a.Host
}

func (entry Entry) Caddyfile() string {
	tls := "# tls disabled"
	a, _ := httpcaddyfile.ParseAddress(entry.Alias)
	// If no scheme is given, default to https.
	if a.Scheme == "" {
		a.Scheme = "https"
	}
	if a.Scheme == "https" {
		tls = strings.TrimSpace(`
	tls {
		issuer internal {
			ca local
			# allow on_demand issuing, but doesn't turn it on
			on_demand
		}
		# turn on on_demand issuing to automatically renew certs
		# when they expire
		on_demand
	}
`)
	}
	return fmt.Sprintf(strings.TrimSpace(`
%s {
	reverse_proxy localhost:%d
	%s
}
	`), entry.Alias, entry.Port, tls)
}
