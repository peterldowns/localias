package hostctl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/exp/slices"
)

type File struct {
	Path  string
	Lines []*Line
}

func (f *File) Add(e *Entry) (*Line, error) {
	line := &Line{Entry: e}
	f.Lines = append(f.Lines, line)
	return line, nil
}

func (f *File) Remove(aliases ...string) ([]*Line, error) {
	for i, alias := range aliases {
		aliases[i] = strings.TrimSpace(alias)
	}
	var kept []*Line
	var removed []*Line
	for _, line := range f.Lines {
		if line.Entry == nil {
			continue
		}
		matched := false
		for _, a := range aliases {
			if slices.Contains(line.Entry.Aliases, a) {
				matched = true
				break
			}
		}
		if matched {
			removed = append(removed, line)
		} else {
			kept = append(kept, line)
		}
	}
	f.Lines = kept
	return removed, nil
}

func (f *File) Contents() string {
	builder := strings.Builder{}
	for _, line := range f.Lines {
		builder.WriteString(line.String())
		builder.WriteString("\n")
	}
	return builder.String()
}

func (f *File) Save(sudo bool) error {
	if f.Path == "" {
		return fmt.Errorf("cannot save file: path is empty")
	}
	var cmd *exec.Cmd
	if sudo {
		cmd = exec.Command("sudo", "tee", f.Path)
	} else {
		cmd = exec.Command("tee", f.Path)
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	go func() {
		defer stdin.Close()
		_, err = io.WriteString(stdin, f.Contents())
		if err != nil {
			panic(err)
		}
	}()
	errtxt, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to save file %s: %w: %s", f.Path, err, errtxt)
	}
	return nil
}

type Line struct {
	Raw   string
	Entry *Entry
}

func (l *Line) String() string {
	if l.Entry != nil && l.Entry.Meta != nil {
		return l.Entry.String()
	}
	return l.Raw
}

type Entry struct {
	IPAddress string
	Aliases   []string
	Meta      *Meta
	Disabled  bool
}

func (e *Entry) String() string {
	l := ""
	if e.Disabled {
		l += "#"
	}
	l += e.IPAddress + "\t"
	l += strings.Join(e.Aliases, "\t")
	if e.Meta != nil {
		b, _ := json.Marshal(e.Meta)
		l += "\t#" + string(b)
	}
	return l
}

type Meta struct {
	Controller string `json:"controller"`
}

func Open(fpath string) (*File, error) {
	fin, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer fin.Close()
	file := Parse(fin)
	file.Path = fpath
	return file, nil
}

func Parse(reader io.Reader) *File {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	var lines []*Line
	for scanner.Scan() {
		line := parseLine(scanner.Text())
		lines = append(lines, line)
	}
	return &File{Path: "", Lines: lines}
}

func parseLine(raw string) *Line {
	entry := parseEntry(raw)
	return &Line{Raw: raw, Entry: entry}
}

func parseEntry(line string) *Entry {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}
	disabled := false
	if isComment(line) {
		disabled = true
		line = line[1:]
	}
	// If it's not a valid line, skip it
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return nil
	}
	ipAddress := fields[0]
	aliases := []string{fields[1]}
	for _, ftxt := range fields[2:] {
		if isComment(ftxt) {
			break
		}
		aliases = append(aliases, ftxt)
	}
	meta := getMeta(line)
	return &Entry{
		IPAddress: ipAddress,
		Aliases:   aliases,
		Meta:      meta,
		Disabled:  disabled,
	}
}

func isComment(s string) bool {
	return strings.HasPrefix(s, "#")
}

func getMeta(line string) *Meta {
	_, comment, found := strings.Cut(line, "#")
	if !found {
		return nil
	}
	meta := Meta{}
	if err := json.Unmarshal([]byte(comment), &meta); err != nil {
		return nil
	}
	return &meta
}
