package hostctl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

const SigilController = "pfpro"

type Line struct {
	Raw   string
	Entry *Entry
}

func (l *Line) String() string {
	if l.Controlled() {
		return l.Entry.String()
	}
	return l.Raw
}

func (l *Line) Controlled() bool {
	return l.Entry != nil && l.Entry.Meta != nil && l.Entry.Meta.Controller == SigilController
}

type Entry struct {
	IPAddress string
	Aliases   []string
	Meta      *Meta
}

func (e *Entry) String() string {
	l := ""
	l += e.IPAddress + "\t"
	l += strings.Join(e.Aliases, "\t")
	if e.Meta != nil {
		b, _ := json.Marshal(e.Meta)
		l += "\t#" + string(b)
	}
	return l
}

type File struct {
	Path  string
	Lines []*Line
}

type Meta struct {
	Controller string `json:"controller"`
}

func (f *File) Add(_ Entry) error {
	return fmt.Errorf("not yet implemented")
}

func (f *File) Remove(_ Entry) error {
	return fmt.Errorf("not yet implemented")
}

func (f *File) Contents() string {
	builder := strings.Builder{}
	for _, line := range f.Lines {
		builder.WriteString(line.String())
		builder.WriteString("\n")
	}
	return builder.String()
}

func Open(fpath string) (*File, error) {
	fin, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer fin.Close()
	file, err := Parse(fin)
	if err != nil {
		return nil, err
	}
	file.Path = fpath
	return file, nil
}

func Parse(reader io.Reader) (*File, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	var lines []*Line
	for scanner.Scan() {
		line, err := parseLine(scanner.Text())
		if err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}
	return &File{Path: "", Lines: lines}, nil
}

func parseLine(raw string) (*Line, error) {
	entry, err := parseEntry(raw)
	if err != nil {
		return nil, err
	}
	return &Line{Raw: raw, Entry: entry}, nil
}

func parseEntry(line string) (*Entry, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, nil
	}
	if isComment(line) {
		return nil, nil
	}
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid line: '%s'", line)
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
	}, nil
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
