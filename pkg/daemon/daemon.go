package daemon

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"

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

// close is used to communicate a shutdown request from the RPC server to the
// proxy server, both of which are running at the same time.
var close chan struct{} //nolint:gochecknoglobals

// Start will fork and start the daemon process, which will start the caddy
// server and the mdns server (if needed) to proxy routes based on the current
// configuration. It also starts an RPC server that will handle requests to
// reload the server with a new configuration or stop the servers and exit the
// daemon.
func Start(cfg *config.Config) error {
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
	if err := server.Start(cfg); err != nil {
		return err
	}

	// Start the daemon RPC server to handle Reload() and Stop() requests from
	// the CLI.
	//
	// See https://eli.thegreenplace.net/2019/unix-domain-sockets-in-go/
	rpcd := new(RPC)
	if err := rpc.Register(rpcd); err != nil {
		panic(err)
	}
	rpc.HandleHTTP()
	sockFile, err := rpcSocketPath()
	if err != nil {
		return err
	}
	if err := os.RemoveAll(sockFile); err != nil {
		return err
	}
	listen, err := net.Listen("unix", sockFile)
	if err != nil {
		return err
	}
	listen.(*net.UnixListener).SetUnlinkOnClose(true)
	defer cleanup(listen.Close)
	srv := &http.Server{}
	close = make(chan struct{}, 1)
	// the RPC Stop() method will send a message on `close`;
	// when that happens, shut down the RPC server and also
	// stop the Caddy / mdns servers.
	go func() {
		<-close
		_ = server.Stop()
		ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
		defer cancel()
		_ = srv.Shutdown(ctx)
		listen.Close()
	}()

	fmt.Println("daemon: listening for rpc requests")
	// Infinite loop to handle RPC requests. The caddy and mdns servers are running
	// in the background.
	if err := srv.Serve(listen); err != nil {
		// calling srv.Shutdown() above causes this function to return with http.ErrServerClosed
		if err != http.ErrServerClosed {
			fmt.Printf("http.Serve exited with err: %s\n", err)
		}
	}
	fmt.Println("daemon: shutdown complete")
	return nil
}

// This struct is used to register RPC commands via the `rpc` package.  This
// lets us send specific commands to a running daemon.
type RPC struct{}

// Reload will stop the existing server (both caddy and mdns) and then start
// them up again with the new config.
func (d RPC) Reload(cfg *config.Config, _ *string) error {
	return server.Start(cfg)
}

// Stop will use the global `close` channel to signal the existing server (both
// caddy and mdns) to shut down cleanly.
func (d RPC) Stop(*string, *string) error {
	close <- struct{}{}
	return nil
}

// cleanup is used to handle defer'd statements that may error.
func cleanup(f func() error) {
	if err := f(); err != nil {
		fmt.Println(err)
	}
}
