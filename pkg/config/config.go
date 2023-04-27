package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Integralist/go-findroot/find"
	"github.com/adrg/xdg"
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/go-yaml/yaml"
)

// Config file is taken from the following list of options, first one that is
// valid is used.
//
// localias --config <path> ...
// LOCALIAS_CONFIGFILE=<path> localias ...
// .localias.yaml in current directory
// .localias.yaml in repository root and command is run from inside a repository
// $XDG_CONFIG_HOME/localias.yaml (or OS fallback if XDG_CONFIG_HOME is not set)
func Path(cfgPath *string) (string, error) {
	if cfgPath != nil && *cfgPath != "" {
		return *cfgPath, nil
	}
	if path := os.Getenv("LOCALIAS_CONFIGFILE"); path != "" {
		return path, nil
	}
	if path := lookup("./.localias.yaml"); path != "" {
		return path, nil
	}
	if repo, err := find.Repo(); err == nil {
		if path := lookup(filepath.Join(repo.Path, ".localias.yaml")); path != "" {
			return path, nil
		}
	}
	return xdg.ConfigFile("localias.yaml")
}

func lookup(path string) string {
	path, err := filepath.Abs(path)
	if err != nil {
		return ""
	}
	if _, err := os.Stat(path); err != nil {
		return ""
	}
	return path
}

func Load(cfgPath *string) (*Config, error) {
	path, err := Path(cfgPath)
	if err != nil {
		return nil, err
	}
	return Open(path)
}

func Open(path string) (*Config, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	contents, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	// Use a MapSlice in order to preserve order
	var entries yaml.MapSlice
	if err := yaml.Unmarshal(contents, &entries); err != nil {
		return nil, err
	}
	c := Config{Path: path}
	for _, entry := range entries {
		c.Upsert(Entry{
			Alias: entry.Key.(string),
			Port:  entry.Value.(int),
		})
	}
	return &c, nil
}

// Upsert will add or update the existing list of entries.  If there is
// already a entry with the same upstream/alias, its downstream/port will be
// updated. If there is not already a entry with the same upstream/alias,
// the new entry will be added to the list.
//
// Returns `true` if an existing entry was updated, `false` if the entry
// was added.
func (c *Config) Upsert(d Entry) bool {
	for i, existing := range c.Entries {
		if existing.Alias == d.Alias {
			c.Entries[i] = d
			return true
		}
	}
	c.Entries = append(c.Entries, d)
	return false
}

func (c *Config) Clear() []Entry {
	removed := c.Entries
	c.Entries = []Entry{}
	return removed
}

func (c *Config) Remove(upstreams ...string) []Entry {
	var removed []Entry
	previous := c.Entries
	c.Clear()
	for _, d := range previous {
		shouldRemove := false
		for _, upstream := range upstreams {
			if d.Alias == upstream {
				shouldRemove = true
				break
			}
		}
		if shouldRemove {
			removed = append(removed, d)
		} else {
			c.Upsert(d)
		}
	}
	return removed
}

type Config struct {
	Path    string
	Entries []Entry
}

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
#   https://secure.test: 9000
#   http://insecure.test: 9001
#   insecure2.test: 9002
#   bareTLD: 9003
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

func (c Config) Caddyfile() string {
	path, err := xdg.StateFile("localias/caddy")
	if err != nil {
		panic(err)
	}
	path, err = filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	// TODO: take an admin port/interface as part of the config settings, and also
	// as part of the CLI?
	global := fmt.Sprintf(strings.TrimSpace(`
{
	admin localhost:2019
	persist_config off
	local_certs
	ocsp_stapling off
	storage file_system "%s"
	pki {
		ca local {
			name localias
			root_cn localias
			intermediate_cn localias
		}
	}
}
`), path)
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

type Entry struct {
	Alias string
	Port  int
}

func (entry Entry) String() string {
	return fmt.Sprintf("%s: %d", entry.Alias, entry.Port)
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
			on_demand
		}
	}
`)
	}
	return fmt.Sprintf(strings.TrimSpace(`
%s {
	reverse_proxy :%d
	%s
}
	`), entry.Alias, entry.Port, tls)
}
