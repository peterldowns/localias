package pfpro

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/go-yaml/yaml"
)

func Load(cfgPath *string) (*Config, error) {
	if cfgPath != nil {
		return ReadConfig(*cfgPath)
	}
	defaultPath, err := DefaultConfigPath()
	if err != nil {
		return nil, err
	}
	return ReadConfig(defaultPath)
}

func DefaultConfigPath() (string, error) {
	path, err := xdg.ConfigFile("pfpro/config.yaml")
	if err != nil {
		return "", err
	}
	return filepath.Abs(path)
}

func WriteConfig(c *Config) error {
	entries := map[string]string{} // TODO: preserve order
	for _, d := range c.Directives {
		entries[d.Upstream] = d.Downstream
	}
	bytes := []byte(strings.TrimSpace(`
# pfpro config file syntax
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

func ReadConfig(path string) (*Config, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	contents, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var entries map[string]string // TODO: preserve order
	if err := yaml.Unmarshal(contents, &entries); err != nil {
		return nil, err
	}
	c := Config{Path: path}
	for upstream, downstream := range entries {
		c.Directives = append(c.Directives, Directive{
			Upstream:   upstream,
			Downstream: downstream,
		})
	}
	return &c, nil
}

type Config struct {
	Path       string
	Directives []Directive
}

type Directive struct {
	Upstream   string
	Downstream string
}

func (directive Directive) String() string {
	return fmt.Sprintf("%s: %s", directive.Upstream, directive.Downstream)
}

func (directive Directive) Caddyfile() string {
	tls := "# tls disabled"
	a, _ := httpcaddyfile.ParseAddress(directive.Upstream)
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
	reverse_proxy %s
	%s
}
	`), directive.Upstream, directive.Downstream, tls)
}

func (c Config) Caddyfile() string {
	path, _ := xdg.ConfigFile("pfpro/caddy/")
	path, _ = filepath.Abs(path)
	global := fmt.Sprintf(strings.TrimSpace(`
{
	admin off
	persist_config off
	local_certs
	ocsp_stapling off
	storage file_system "%s"
	pki {
		ca local {
			name pfpro
			root_cn pfpro
			intermediate_cn pfpro
		}
	}
}
`), path)
	blocks := []string{global}
	for _, x := range c.Directives {
		blocks = append(blocks, x.Caddyfile())
	}
	// extra newline prevents "caddy fmt" warning in logs
	return strings.Join(blocks, "\n") + "\n"
}
