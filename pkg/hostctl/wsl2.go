package hostctl

import "os"

type WSL2Controller struct {
	EtcHostsController *FileController
	TmpController      *FileController
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

	tmpfile, err := os.CreateTemp("", "localias-hosts-*")
	if err != nil {
		panic(err)
	}
	defer tmpfile.Close()

	// TODO: get the contents from the WSL hosts file
	if _, err := tmpfile.Write([]byte("127.0.0.1 localhost")); err != nil {
		panic(err)
	}

	// defer os.Remove(f.Name()) // clean up
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

func (w *WSL2Controller) Apply() error {
	if err := w.EtcHostsController.Apply(); err != nil {
		return err
	}
	if err := w.TmpController.Apply(); err != nil {
		return err
	}
	// TODO: now sync back to the windows host
	return nil
}
