package hostctl

import (
	"os"

	"github.com/peterldowns/localias/pkg/wsl"
)

type WSL2Controller struct {
	EtcHostsController *FileController
	TmpController      *FileController
	IP                 string
}

func NewWSL2Controller() *WSL2Controller {
	x := WSL2Controller{}
	name := "localias"
	dryrun := false
	x.EtcHostsController = &FileController{
		HostsFile: "/etc/hosts",
		Sudo:      true,
		DryRun:    dryrun,
		Name:      name,
	}
	x.IP = wsl.IP()
	tmpfile, err := os.CreateTemp("", "localias-hosts-*")
	if err != nil {
		panic(err)
	}
	defer tmpfile.Close()

	contents, err := wsl.ReadWindowsHosts()
	if err != nil {
		panic(err)
	}
	if _, err := tmpfile.Write([]byte(contents)); err != nil {
		panic(err)
	}
	x.TmpController = &FileController{
		HostsFile: tmpfile.Name(),
		Sudo:      false,
		DryRun:    dryrun,
		Name:      name,
	}
	return &x
}

func (w *WSL2Controller) Set(ip string, alias string) error {
	if err := w.EtcHostsController.Set(ip, alias); err != nil {
		return err
	}
	if err := w.TmpController.Set(ip, alias); err != nil {
		return err
	}
	return nil
}

func (w *WSL2Controller) SetLocal(alias string) error {
	if err := w.EtcHostsController.SetLocal(alias); err != nil {
		return err
	}
	if err := w.TmpController.Set(w.IP, alias); err != nil {
		return err
	}
	return nil
}

func (w *WSL2Controller) Remove(alias string) error {
	if err := w.EtcHostsController.Remove(alias); err != nil {
		return err
	}
	if err := w.TmpController.Remove(alias); err != nil {
		return err
	}
	return nil
}

func (w *WSL2Controller) List() ([]*Line, error) {
	return w.TmpController.List()
}

func (w *WSL2Controller) Clear() error {
	if err := w.EtcHostsController.Clear(); err != nil {
		return err
	}
	if err := w.TmpController.Clear(); err != nil {
		return err
	}
	return nil
}

// TODO: this is definitely the wrong abstraction for dealing wth
// multiple controllers. fuck it.
func (w *WSL2Controller) Apply() (bool, error) {
	if _, err := w.EtcHostsController.Apply(); err != nil {
		return false, err
	}
	changes, err := w.TmpController.Apply()
	if err != nil {
		return changes, err
	}
	if changes {
		return true, wsl.WriteWindowsHostsFromFile(w.TmpController.HostsFile)
	}
	return false, nil
}
