package hostctl

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"strings"
)

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

func Open(fpath string) ([]*Line, error) {
	fin, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer fin.Close()
	return Parse(fin), nil
}

func Parse(reader io.Reader) []*Line {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	var lines []*Line
	for scanner.Scan() {
		line := parseLine(scanner.Text())
		lines = append(lines, line)
	}
	return lines
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
	// Controlled lines only have 1 alias
	aliases := []string{fields[1]}
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
