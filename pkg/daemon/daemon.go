package daemon

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/adrg/xdg"
	caddycmd "github.com/caddyserver/caddy/v2/cmd"
	godaemon "github.com/sevlyar/go-daemon"

	"github.com/peterldowns/localias/pkg/config"
	"github.com/peterldowns/localias/pkg/hostctl"
	"github.com/peterldowns/localias/pkg/server"
)

// Start will apply the latest configuration and start the caddy daemon server,
// then exit. If the caddy daemon server is already running, it will exit with
// an error.
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
	if err := server.Start(hctl, cfg); err != nil {
		return err
	}
	select {} //nolint:revive // valid empty block, keeps the server running forever.
}

// Status will determine whether or not the caddy daemon server is running.  If
// it is, it returns the non-nil os.Process of that daemon.
func Status() (*os.Process, error) {
	cntxt := daemonContext()
	proc, err := cntxt.Search()
	if err != nil {
		// If the pidfile for the daemon doesn't exist, cntxt.Search() throws a
		// PathError. In that case, we assume the daemon is not running, and
		// return nil.
		var pathError *os.PathError
		if errors.As(err, &pathError) {
			return nil, nil
		}
	}
	return proc, nil
}

// Stop will attempt to stop the caddy daemon server by sending an API request
// over http. If the daemon server is not running, it will return an error.
func Stop(cfg *config.Config) error {
	address, err := determineAPIAddress(cfg)
	if err != nil {
		return fmt.Errorf("could not determine api address: %w", err)
	}
	resp, err := caddycmd.AdminAPIRequest(address, http.MethodPost, "/stop", nil, nil)
	if err != nil {
		return fmt.Errorf("request to /stop failed: %w", err)
	}
	defer resp.Body.Close()
	return nil
}

// Reload will apply the latest configuration details, and then update the
// running caddy daemon server's configuration by sending an API request over
// http.  If the daemon server is not running, it will return an error.
func Reload(hctl *hostctl.Controller, cfg *config.Config) error {
	err := config.Apply(hctl, cfg)
	if err != nil {
		return err
	}
	cfgJSON, _, err := cfg.CaddyJSON()
	if err != nil {
		return err
	}
	address, err := determineAPIAddress(cfg)
	if err != nil {
		return err
	}
	headers := make(http.Header)
	headers.Set("Cache-Control", "must-revalidate")
	resp, err := caddycmd.AdminAPIRequest(address, http.MethodPost, "/load", headers, bytes.NewReader(cfgJSON))
	if err != nil {
		return fmt.Errorf("failed to send config to daemon: %w", err)
	}
	defer resp.Body.Close()
	return nil
}

// daemonContext returns a consistent godaemon context that is used to control
// the caddy daemon server.
func daemonContext() *godaemon.Context {
	pidFile, err := xdg.StateFile("localias/daemon.pid")
	if err != nil {
		panic(err)
	}
	return &godaemon.Context{
		PidFileName: pidFile,
		PidFilePerm: 0o644,
		Umask:       0o27,
	}
}
