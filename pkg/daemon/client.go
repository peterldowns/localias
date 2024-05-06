package daemon

import (
	"errors"
	"os"
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
