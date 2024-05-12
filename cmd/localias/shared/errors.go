package shared

import (
	"fmt"
	"strings"
)

func ConvertErr(err error) error {
	if err == nil {
		return nil
	}
	errMsg := err.Error()
	if strings.Contains(errMsg, "address already in use") {
		return DaemonRunningError{}
	}
	if strings.Contains(errMsg, "connect: connection refused") {
		return DaemonNotRunningError{}
	}
	if strings.Contains(errMsg, "bind: permission denied") {
		return BindNotAllowedError{}
	}
	return err
}

type LocaliasError interface {
	Error() string
	Code() string
}

// bind: permission denied
type BindNotAllowedError struct{}

func (BindNotAllowedError) Error() string {
	return "the server is not allowed to bind to ports 443/80"
}

func (BindNotAllowedError) Code() string {
	return "privports"
}

type DaemonNotRunningError struct{}

func (DaemonNotRunningError) Error() string {
	return "the localias daemon is not running"
}

func (DaemonNotRunningError) Code() string {
	return "daemon_not_running"
}

type DaemonRunningError struct {
	Pid int
}

func (x DaemonRunningError) Error() string {
	if x.Pid != 0 {
		return fmt.Sprintf("the localias daemon is already running (pid=%d)", x.Pid)
	}
	return strings.TrimSpace(`
localias could not start successfully. Most likely there is another instance of
localias or some other kind of proxy or server listening to ports 443/80, which
is preventing another instance from starting. Common causes:

- You have another instance of localias running in a different terminal
- You have a proxy server like Caddy, Nginx, or Apache running
- There is a bug in localias

Please see the https://github.com/peterldowns/localias README for some
diagnostics and ideas for how to debug this.
`)
}

func (DaemonRunningError) Code() string {
	return "daemon_running"
}
