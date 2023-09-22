package daemon

import (
	"errors"
	"net/rpc"
	"os"

	"github.com/adrg/xdg"

	"github.com/peterldowns/localias/pkg/config"
)

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

// Reload requests that the running daemon reload and use the given config.
func Reload(cfg *config.Config) error {
	var ignored string // necessary for RPC interface
	return client().Call("RPC.Reload", cfg, &ignored)
}

// Stop requests that the running daemon exit cleanly and shut down.
func Stop() error {
	// Send a request to the daemon to shut down.
	var ignored string // necessary for RPC interface
	return client().Call("RPC.Stop", &ignored, &ignored)
}

// Kill force-kills the daemon if it's running.
func Kill() error {
	proc, err := daemonContext().Search()
	if err != nil {
		return err
	}
	if proc != nil {
		return proc.Kill()
	}
	return nil
}

func client() *rpc.Client {
	sockFile, err := rpcSocketPath()
	if err != nil {
		panic(err)
	}
	client, err := rpc.DialHTTP("unix", sockFile)
	if err != nil {
		panic(err)
	}
	return client
}

func rpcSocketPath() (string, error) {
	return xdg.StateFile("localias/daemon.sock")
}
