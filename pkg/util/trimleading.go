package util

import (
	"strings"

	"github.com/fatih/color"
)

// Example is a helper for generating CLI example docs for cobra commands.  It
// removes any surrounding space from a string, then removes any leading
// whitespace from each line in the string. Any comments in the string will be
// colored as faint.
func Example(s string) string {
	in := strings.Split(strings.TrimSpace(s), "\n")
	var out []string

	for _, x := range in {
		x = strings.TrimSpace(x)
		if len(x) > 0 && x[0] == '#' {
			x = color.New(color.Faint).Sprint(x)
		}
		out = append(out, "  "+x)
	}
	return strings.Join(out, "\n")
}
