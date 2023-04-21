package daemon

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	caddycmd "github.com/caddyserver/caddy/v2/cmd"
	godaemon "github.com/sevlyar/go-daemon"

	"github.com/peterldowns/localias/pkg/config"
	"github.com/peterldowns/localias/pkg/hostctl"
)

// TODO: take some details here from the config,
func daemonContext() *godaemon.Context {
	return &godaemon.Context{
		PidFileName: "localias.pid",
		PidFilePerm: 0o644,
		WorkDir:     "./",
		Umask:       0o27,
	}
}

func Start(hctl *hostctl.Controller, cfg *config.Config) error {
	existing, err := Status()
	if err != nil {
		return err
	}
	if existing != nil {
		return fmt.Errorf("daemon is already running")
	}

	cntxt := daemonContext()
	d, err := cntxt.Reborn()
	if err != nil {
		return err
	}
	// parent process: the child has started, so exit
	if d != nil {
		return nil
	}
	// child process: defer a cleanup function, then run
	defer func() {
		_ = cntxt.Release()
	}()
	return Run(hctl, cfg)
}

func Status() (*os.Process, error) {
	cntxt := daemonContext()
	return cntxt.Search()
}

// TODO: check signature of this fn
func Stop(_ *hostctl.Controller, cfg *config.Config) error {
	// TODO: use admin API to stop the daemon
	caddyfile := cfg.Caddyfile()
	cfgAdapter := caddyconfig.GetAdapter("caddyfile")
	if cfgAdapter == nil {
		panic(fmt.Errorf("nil cfgadapter"))
	}
	xAddr := "" // TODO: maybe allow passing this?
	config, _, err := cfgAdapter.Adapt([]byte(caddyfile), map[string]any{
		"filename": "Caddyfile",
	})
	if err != nil {
		fmt.Printf("failed to get cfg adapter i think: %v\n", err)
		return err
	}
	// https://github.com/caddyserver/caddy/blob/be53e432fcac0a9b9accbc36885304639e8ca70b/cmd/commandfuncs.go#L146
	configFile := ""
	configFileAdapter := ""
	address, err := caddycmd.DetermineAdminAPIAddress(xAddr, config, configFile, configFileAdapter)
	if err != nil {
		return fmt.Errorf("could not determine api address: %w", err)
	}
	resp, err := caddycmd.AdminAPIRequest(address, http.MethodPost, "/stop", nil, nil)
	if err != nil {
		return fmt.Errorf("failed to stop daemon: %w", err)
	}
	defer resp.Body.Close()
	return nil
}

// TODO: check signature of this fn
func Reload(hctl *hostctl.Controller, cfg *config.Config) error {
	if err := hctl.Clear(); err != nil {
		return err
	}
	for _, directive := range cfg.Directives {
		up, err := httpcaddyfile.ParseAddress(directive.Upstream)
		if err != nil {
			return err
		}
		if err := hctl.Set("127.0.0.1", up.Host); err != nil {
			return err
		}
	}
	if err := hctl.Apply(); err != nil {
		return err
	}

	caddyfile := cfg.Caddyfile()
	cfgAdapter := caddyconfig.GetAdapter("caddyfile")
	if cfgAdapter == nil {
		panic(fmt.Errorf("nil cfgadapter"))
	}
	config, warnings, err := cfgAdapter.Adapt([]byte(caddyfile), map[string]any{
		"filename": "Caddyfile",
	})
	if err != nil {
		fmt.Printf("failed to get cfg adapter i think: %v\n", err)
		return err
	}
	for _, w := range warnings {
		fmt.Printf("[warning]: %s\n", w.String())
	}
	// https://github.com/caddyserver/caddy/blob/be53e432fcac0a9b9accbc36885304639e8ca70b/cmd/commandfuncs.go#L146
	configFile := ""
	configFileAdapter := ""
	xAddr := "" // TODO: maybe allow passing this?
	address, err := caddycmd.DetermineAdminAPIAddress(xAddr, config, configFile, configFileAdapter)
	if err != nil {
		return fmt.Errorf("could not determine api address: %w", err)
	}
	headers := make(http.Header)
	headers.Set("Cache-Control", "must-revalidate")
	resp, err := caddycmd.AdminAPIRequest(address, http.MethodPost, "/load", headers, bytes.NewReader(config))
	if err != nil {
		return fmt.Errorf("failed to send config to daemon: %w", err)
	}
	defer resp.Body.Close()
	return nil
}

func Run(hctl *hostctl.Controller, cfg *config.Config) error {
	if err := hctl.Clear(); err != nil {
		return err
	}
	for _, directive := range cfg.Directives {
		up, err := httpcaddyfile.ParseAddress(directive.Upstream)
		if err != nil {
			return err
		}
		if err := hctl.Set("127.0.0.1", up.Host); err != nil {
			return err
		}
	}
	if err := hctl.Apply(); err != nil {
		return err
	}

	caddyfile := cfg.Caddyfile()
	cfgAdapter := caddyconfig.GetAdapter("caddyfile")
	if cfgAdapter == nil {
		panic(fmt.Errorf("nil cfgadapter"))
	}
	config, warnings, err := cfgAdapter.Adapt([]byte(caddyfile), map[string]any{
		"filename": "Caddyfile",
	})
	if err != nil {
		fmt.Printf("failed to get cfg adapter i think: %v\n", err)
		return err
	}
	for _, w := range warnings {
		fmt.Printf("[warning]: %s\n", w.String())
	}
	// TODO: actually daemonize like
	// https://github.com/caddyserver/caddy/blob/be53e432fcac0a9b9accbc36885304639e8ca70b/cmd/commandfuncs.go#L146
	// reloading should be as simple as calling caddy.Load() again
	err = caddy.Load(config, false)
	if err != nil {
		fmt.Printf("failing with %v\n", err)
		return err
	}
	select {} //nolint:revive // valid empty block
}
