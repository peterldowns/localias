package hostctl

import "fmt"

const (
	DefaultHostsFile = "/etc/hosts"
	DefaultSudo      = true
	DefaultDryRun    = false
	DefaultName      = "pfpro"
)

var ErrFileNotOpen = fmt.Errorf("file is nil, call .Open() first")

type Controller struct {
	HostsFile string
	Sudo      bool
	DryRun    bool
	Name      string
	file      *File
}

func DefaultController() *Controller {
	return NewController(
		DefaultHostsFile,
		DefaultSudo,
		DefaultDryRun,
		DefaultName,
	)
}

func NewController(
	hostsFile string,
	sudo bool,
	dryRun bool,
	name string,
) *Controller {
	return &Controller{
		HostsFile: hostsFile,
		Sudo:      sudo,
		DryRun:    dryRun,
		Name:      name,
	}
}

func (c *Controller) Open() (*File, error) {
	return Open(c.HostsFile)
}

func (c *Controller) Read() error {
	f, err := Open(c.HostsFile)
	if err != nil {
		return err
	}
	c.file = f
	return nil
}

func (c *Controller) Save() error {
	if c.DryRun {
		return nil
	}
	file, err := c.File()
	if err != nil {
		return err
	}
	return file.Save(c.Sudo)
}

func (c *Controller) File() (*File, error) {
	if c.file == nil {
		if err := c.Read(); err != nil {
			return nil, err
		}
	}
	return c.file, nil
}

func (c *Controller) Add(ip string, aliases ...string) ([]*Line, error) {
	file, err := c.File()
	if err != nil {
		return nil, err
	}
	var added []*Line
	for _, alias := range aliases {
		line, err := file.Add(&Entry{
			IPAddress: ip,
			Aliases:   []string{alias},
			Meta: &Meta{
				Controller: c.Name,
			},
		})
		if err != nil {
			return nil, err
		}
		added = append(added, line)
	}
	return added, nil
}

func (c *Controller) List() ([]*Line, error) {
	file, err := c.File()
	if err != nil {
		return nil, err
	}
	var lines []*Line
	for _, line := range file.Lines {
		if !controlled(line) {
			continue
		}
		lines = append(lines, line)
	}
	return lines, nil
}

func controlled(l *Line) bool {
	if l.Entry == nil {
		return false
	}
	if l.Entry.Meta == nil {
		return false
	}
	return l.Entry.Meta.Controller == DefaultName
}

func (c *Controller) Remove(aliases ...string) ([]*Line, error) {
	f, err := Open(c.HostsFile)
	if err != nil {
		return nil, err
	}
	removed, err := f.Remove(aliases...)
	if err != nil {
		return nil, err
	}
	if c.DryRun {
		return removed, nil
	}
	if err := f.Save(c.Sudo); err != nil {
		return nil, err
	}
	return removed, nil
}
