package daemon

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/adrg/xdg"
	godaemon "github.com/sevlyar/go-daemon"

	"github.com/peterldowns/localias/pkg/config"
	"github.com/peterldowns/localias/pkg/server"
)

// daemonContext returns a consistent go-daemon context that is used to control
// the daemon process.
func daemonContext() *godaemon.Context {
	logFile, _ := xdg.StateFile("localias/daemon.log")
	pidFile, err := xdg.StateFile("localias/daemon.pid")
	if err != nil {
		panic(err)
	}
	return &godaemon.Context{
		LogFileName: logFile,
		LogFilePerm: 0o644,
		PidFileName: pidFile,
		PidFilePerm: 0o644,
		Umask:       0o27,
	}
}

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
	if err := instance.StartCaddy(); err != nil {
		return err
	}
	fmt.Println("caddy started successfully")
	cntxt := daemonContext()
	proc, err := cntxt.Reborn()
	if err != nil {
		return err
	}
	// parent process, after fork succeeds
	if proc != nil {
		return nil
	}
	// child (daemon) process, after fork succeeds
	defer cleanup(cntxt.Release) // removes the go-daemon PID file on shutdown.
	return Run(cfg)
}

// run is the logic that the daemon runs after it is forked.  It will loop until
// it receives a shutdown request, after which the function exits cleanly.
func Run(cfg *config.Config) error {
	// Start the caddy proxy server and the mdns responders.
	fmt.Print("daemon: running!")
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
	<-quitChannel
	return nil
}

// cleanup is used to handle defer'd statements that may error.
func cleanup(f func() error) {
	if err := f(); err != nil {
		fmt.Println(err)
	}
}
