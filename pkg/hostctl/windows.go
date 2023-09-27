package hostctl

import (
	"fmt"
	"os"
)

type WindowsController struct {
	FileController
	WindowsHostsFile string
}

// NewWindowsController creates a new [WindowsController]. It can be used when
// localias is ran on Windows (GOOS=windows). Contents for the hosts file are first
// written to a temporary file and then compared to contain changes. If there are
// changes, the actual Windows hosts file is modified. This is the same model the
// WSLController uses, except for WSL specific PowerShell script invocations.
func NewWindowsController(name string) *WindowsController {
	x := WindowsController{}
	x.WindowsHostsFile = os.Getenv("SystemRoot") + `\System32\drivers\etc\hosts`

	// TODO: don't panic in this code, but return an error instead? But it's
	// the pattern that's in use, and seems to be used for error handling too.
	tmpfile, err := os.CreateTemp("", fmt.Sprintf("%s-hosts-*", name))
	if err != nil {
		panic(err)
	}
	defer tmpfile.Close()

	contents, err := os.ReadFile(x.WindowsHostsFile)
	if err != nil {
		panic(err)
	}
	if _, err := tmpfile.Write([]byte(contents)); err != nil {
		panic(err)
	}

	x.FileController = FileController{
		Path: tmpfile.Name(),
		Sudo: true, // always needs administrator permissions
		Name: name,
	}

	return &x
}

func (w *WindowsController) Set(ip string, alias string) error {
	return w.FileController.Set(ip, alias)
}

func (w *WindowsController) SetLocal(alias string) error {
	return w.FileController.Set("127.0.0.1", alias)
}

func (w *WindowsController) Remove(alias string) error {
	return w.FileController.Remove(alias)
}

func (w *WindowsController) Clear() error {
	return w.FileController.Clear()
}

func (w *WindowsController) Apply() (bool, error) {
	changes, err := w.FileController.Apply()
	if err != nil {
		return false, err
	}
	if changes {
		// read the current state of the temporary file after changes
		// were detected.
		contents, err := os.ReadFile(w.FileController.Path)
		if err != nil {
			return false, err
		}

		// TODO(hs): use PowerShell instead, similar to the WSL implementation, so that similar errors
		// are reported and can be handled in a unified way?
		if err := os.WriteFile(w.WindowsHostsFile, contents, 0o644); err != nil {
			return false, err
		}

		return true, nil
	}
	return false, nil
}

func (w *WindowsController) List() (map[string][]*Line, error) {
	return w.FileController.List()
}
