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

	"github.com/peterldowns/localias/cmd/localias/shared"
	"github.com/peterldowns/localias/pkg/config"
	"github.com/peterldowns/localias/pkg/server"
)

// Start will apply the latest configuration and then start the daemon process.
// If it is already running, this function will return an error.
func Start(cfg *config.Config) error {
	existing, err := Status()
	if err != nil {
		return err
	}
	if existing != nil {
		return shared.DaemonRunning{Pid: existing.Pid}
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
	if err := server.Start(cfg); err != nil {
		return err
	}
	select {}
}

// Status will determine whether or not the daemon process is running. If
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

// Stop will attempt to stop the daemon process by sending an API request
// over http. If the daemon process is not running, this will return an error.
func Stop(cfg *config.Config) error {
	existing, err := Status()
	if err != nil {
		return err
	}
	if existing == nil {
		return shared.DaemonNotRunning{}
	}
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
// running daemon process's server configuration by sending an API request over
// http. If the daemon process is not running, this will return an error.
func Reload(cfg *config.Config) error {
	existing, err := Status()
	if err != nil {
		return err
	}
	if existing == nil {
		return shared.DaemonNotRunning{}
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
// the daemon process that will run the caddy server.
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
