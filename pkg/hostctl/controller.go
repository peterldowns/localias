package hostctl

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

const (
	DefaultHostsFile = "/etc/hosts"
	DefaultSudo      = true
	DefaultDryRun    = false
	DefaultName      = "localias"
)

var ErrFileNotOpen = fmt.Errorf("file is nil, call .Open() first")

type Controller interface {
	Set(ip string, alias string) error
	Remove(alias string) error
	List() ([]*Line, error)
	Clear() error
	Apply() error
}

type FileController struct {
	HostsFile string
	Sudo      bool
	DryRun    bool
	Name      string
	// Internal details
	lines []*Line
	lmap  map[string]string
}

func NewController(
	hostsFile string,
	sudo bool,
	dryRun bool,
	name string,
) *FileController {
	return &FileController{
		HostsFile: hostsFile,
		Sudo:      sudo,
		DryRun:    dryRun,
		Name:      name,
	}
}

func DefaultController() *FileController {
	return NewController(
		DefaultHostsFile,
		DefaultSudo,
		DefaultDryRun,
		DefaultName,
	)
}

func (c *FileController) read() error {
	if c.lines == nil {
		lines, err := Open(c.HostsFile)
		if err != nil {
			return err
		}
		c.lines = lines
	}
	if c.lmap == nil {
		c.lmap = make(map[string]string)
		for _, l := range c.lines {
			if isControlled(l, c.Name) {
				c.lmap[l.Entry.Aliases[0]] = l.Entry.IPAddress
			}
		}
	}
	return nil
}

func (c *FileController) Set(ip string, alias string) error {
	if err := c.read(); err != nil {
		return err
	}
	c.lmap[alias] = ip
	return nil
}

func (c *FileController) Remove(alias string) error {
	if err := c.read(); err != nil {
		return err
	}
	delete(c.lmap, alias)
	return nil
}

func (c *FileController) List() ([]*Line, error) {
	if err := c.read(); err != nil {
		return nil, err
	}
	var toDisplay []*Line
	for _, line := range c.lines {
		if isControlled(line, c.Name) {
			toDisplay = append(toDisplay, line)
		}
	}
	return toDisplay, nil
}

func (c *FileController) Clear() error {
	if err := c.read(); err != nil {
		return err
	}
	c.lmap = make(map[string]string)
	return nil
}

func (c *FileController) Apply() error {
	if err := c.read(); err != nil {
		return err
	}
	var changed bool
	var result []*Line
	for _, line := range c.lines {
		if !isControlled(line, c.Name) {
			result = append(result, line)
			continue
		}
		alias := line.Entry.Aliases[0]
		if ip, ok := c.lmap[alias]; ok {
			if line.Entry.IPAddress != ip {
				line.Entry.IPAddress = ip
				changed = true // modified an existing line
			}
			result = append(result, line)
			delete(c.lmap, alias)
			continue
		}
		changed = true // removed an existing line
	}
	for alias, ip := range c.lmap {
		l := &Line{
			Entry: &Entry{
				IPAddress: ip,
				Aliases:   []string{alias},
				Meta: &Meta{
					Controller: c.Name,
				},
			},
		}
		changed = true // added a new line
		result = append(result, l)
	}
	c.lines = result
	if !changed {
		return nil
	}
	return c.save()
}

func (c *FileController) save() error {
	if c.DryRun {
		return nil
	}
	if c.HostsFile == "" {
		return fmt.Errorf("cannot save file: path is empty")
	}
	var cmd *exec.Cmd
	if c.Sudo {
		cmd = exec.Command("sudo", "tee", c.HostsFile)
	} else {
		cmd = exec.Command("tee", c.HostsFile)
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	go func() {
		defer stdin.Close()
		_, err = io.WriteString(stdin, asFile(c.lines))
		if err != nil {
			panic(err)
		}
	}()
	if errtxt, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to save file %s: %w: %s", c.HostsFile, err, errtxt)
	}
	return nil
}

func asFile(lines []*Line) string {
	builder := strings.Builder{}
	for _, line := range lines {
		builder.WriteString(line.String())
		builder.WriteString("\n")
	}
	return builder.String()
}

func isControlled(line *Line, controllerName string) bool {
	if line.Entry == nil {
		return false
	}
	if line.Entry.Meta == nil {
		return false
	}
	return line.Entry.Meta.Controller == controllerName
}
