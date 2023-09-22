// caddymodules is a package that imports all the plugins necessary to make the
// embedded caddy proxy server work correctly. After importing _this_ package, a
// golang program can call `caddy.Load(...)` to start the proxy server.
//
// Normally a golang program could import
// `"github.com/caddyserver/caddy/v2/modules/standard"` to achieve the same
// goal, but something in the nix/golang module building code chokes on some
// deep dependency in the caddy modules. Through trial and error, this set of
// imports seems to work.
package caddymodules

import (
	// contents of github.com/caddyserver/caddy/v2/modules/standard
	// excluding
	// _ "github.com/caddyserver/caddy/v2/modules/caddyhttp/standard"
	// because of gomod2nix bug
	_ "github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	_ "github.com/caddyserver/caddy/v2/modules/caddyevents"
	_ "github.com/caddyserver/caddy/v2/modules/caddyevents/eventsconfig"
	_ "github.com/caddyserver/caddy/v2/modules/caddypki"
	_ "github.com/caddyserver/caddy/v2/modules/caddypki/acmeserver"
	_ "github.com/caddyserver/caddy/v2/modules/caddytls"
	_ "github.com/caddyserver/caddy/v2/modules/caddytls/distributedstek"
	_ "github.com/caddyserver/caddy/v2/modules/caddytls/standardstek"
	_ "github.com/caddyserver/caddy/v2/modules/filestorage"
	_ "github.com/caddyserver/caddy/v2/modules/logging"
	_ "github.com/caddyserver/caddy/v2/modules/metrics"

	// contents of github.com/caddyserver/caddy/v2/modules/caddyhttp/standard
	// excluding
	// _ "github.com/caddyserver/caddy/v2/modules/caddyhttp/tracing"
	// because of gomod2nix bug
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/caddyauth"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/encode"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/encode/brotli"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/encode/gzip"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/encode/zstd"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/fileserver"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/headers"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/map"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/push"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/requestbody"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/reverseproxy"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/reverseproxy/fastcgi"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/reverseproxy/forwardauth"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/rewrite"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/templates"
)
