package hostctl

import (
	"fmt"
	"os"

	"github.com/peterldowns/localias/pkg/wsl"
)

type WSLController struct {
	FileController
	wslIP string
}

func NewWSLController(name string) *WSLController {
	x := WSLController{}
	x.wslIP = wsl.IP()
	tmpfile, err := os.CreateTemp("", fmt.Sprintf("%s-hosts-*", name))
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
	x.FileController = FileController{
		Path: tmpfile.Name(),
		Sudo: false,
		Name: name,
	}
	return &x
}

func (w *WSLController) Set(ip string, alias string) error {
	return w.FileController.Set(ip, alias)
}

func (w *WSLController) SetLocal(alias string) error {
	return w.FileController.Set(w.wslIP, alias)
}

func (w *WSLController) Remove(alias string) error {
	return w.FileController.Remove(alias)
}

func (w *WSLController) Clear() error {
	return w.FileController.Clear()
}

func (w *WSLController) Apply() (bool, error) {
	changes, err := w.FileController.Apply()
	if err != nil {
		return false, err
	}
	if changes {
		return true, wsl.WriteWindowsHostsFromFile(w.FileController.Path)
	}
	return false, nil
}
