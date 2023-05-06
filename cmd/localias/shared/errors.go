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
		return DaemonRunning{}
	}
	if strings.Contains(errMsg, "connect: connection refused") {
		return DaemonNotRunning{}
	}
	if strings.Contains(errMsg, "bind: permission denied") {
		return BindNotAllowed{}
	}
	return err
}

type LocaliasError interface {
	Error() string
	Code() string
}

// bind: permission denied
type BindNotAllowed struct{}

func (x BindNotAllowed) Error() string {
	return "the server is not allowed to bind to ports 443/80"
}

func (x BindNotAllowed) Code() string {
	return "privports"
}

type DaemonNotRunning struct{}

func (x DaemonNotRunning) Error() string {
	return "the localias daemon is not running"
}

func (x DaemonNotRunning) Code() string {
	return "daemon_not_running"
}

type DaemonRunning struct {
	Pid int
}

func (x DaemonRunning) Error() string {
	if x.Pid != 0 {
		return fmt.Sprintf("the localias daemon is already running (pid=%d)", x.Pid)
	}
	return "a localias daemon or server is already running"
}

func (x DaemonRunning) Code() string {
	return "daemon_running"
}
