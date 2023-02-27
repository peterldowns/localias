package pfpro

import (
	"fmt"
	"strings"

	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
)

type Config []Directive

type Directive struct {
	Upstream   string
	Downstream string
}

func (directive Directive) Caddyfile() string {
	tls := "# tls disabled"
	a, _ := httpcaddyfile.ParseAddress(directive.Upstream)
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
	global := strings.TrimSpace(`
{
	admin off
	persist_config off
	local_certs
	ocsp_stapling off
	storage file_system /Users/pd/.config/pfpro
	pki {
		ca local {
			name pfpro
			root_cn pfpro
			intermediate_cn pfpro
		}
	}
}
`)
	blocks := []string{global}
	for _, x := range c {
		blocks = append(blocks, x.Caddyfile())
	}
	// extra newline prevents "caddy fmt" warning in logs
	return strings.Join(blocks, "\n") + "\n"
}
