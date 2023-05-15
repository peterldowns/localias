package hostctl

import (
	"fmt"
)

var ErrFileNotOpen = fmt.Errorf("file is nil, call .Open() first")

type Controller interface {
	Set(ip string, alias string) error
	SetLocal(alias string) error
	Remove(alias string) error
	Clear() error
	Apply() (bool, error)
	List() (map[string][]*Line, error)
}
