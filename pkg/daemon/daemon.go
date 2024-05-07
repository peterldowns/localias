package daemon

import (
	"errors"
	"os"

	"github.com/adrg/xdg"
	godaemon "github.com/sevlyar/go-daemon"

	"github.com/peterldowns/localias/pkg/config"
	"github.com/peterldowns/localias/pkg/server"
)

// Start will fork a daemon process that loops until it receives
// sigint/sigterm/sigabrt.
// Start will fork and start the daemon process, which will start the caddy
// server and the mdns server (if needed) to proxy routes based on the current
// configuration. It also starts an RPC server that will handle requests to
// reload the server with a new configuration or stop the servers and exit the
// daemon.
func Start(cfg *config.Config) error {
	if err := Kill(); err != nil {
		return err
	}
	instance := &server.Server{Config: cfg}
	if err := instance.Start(); err != nil {
		return err
	}
	cntxt := daemonContext()
	proc, err := cntxt.Reborn()
	if err != nil {
		return err
	}
	// parent process, after fork succeeds just exit
	if proc != nil {
		return nil
	}
	// child (daemon) process, after fork succeeds. the caddy/mdns servers are
	// started and already running; just need to run until
	// sigint/sigterm/sigabrt. waiting on the channel is cheaper than
	// spinlocking and more correct, too.
	server.WaitForExitSignal()
	return cntxt.Release()
}

// Status will determine whether or not the daemon process is running. If it is,
// it returns the non-nil os.Process of that daemon.
func Status() (*os.Process, error) {
	cntxt := daemonContext()
	proc, err := cntxt.Search()
	if err != nil && !isPathError(err) {
		return nil, err
	}
	return proc, nil
}

// Kill force-kills the daemon if it's running, otherwise does nothing.
func Kill() error {
	proc, err := daemonContext().Search()
	if err != nil && !isPathError(err) {
		return err
	}
	if proc != nil {
		return proc.Kill()
	}
	return nil
}

// If the pidfile for the daemon doesn't exist, cntxt.Search() throws a
// PathError. In that case, we assume the daemon is not running, and return nil.
func isPathError(err error) bool {
	var pathError *os.PathError
	return errors.As(err, &pathError)
}

// daemonContext returns a consistent go-daemon context that is used to control
// the daemon process.
func daemonContext() *godaemon.Context {
	logFile, _ := xdg.StateFile("localias/daemon.log")
	pidFile, _ := xdg.StateFile("localias/daemon.pid")
	return &godaemon.Context{
		LogFileName: logFile,
		LogFilePerm: 0o644,
		PidFileName: pidFile,
		PidFilePerm: 0o644,
		Umask:       0o27,
	}
}
