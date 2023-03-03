package hostctl

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseUncontrolledRoundtrips(t *testing.T) {
	contents := etcfile(`
# just a comment     


127.0.0.1	localhost	host.local
`)
	lines := Parse(strings.NewReader(contents))
	require.NotNil(t, lines)
	require.Len(t, lines, 4)
	require.Equal(t, contents, asFile(lines))
}

func TestParseControlled(t *testing.T) {
	contents := etcfile(`
# /etc/hosts file blah blah blah
127.0.0.1 localhost


127.0.0.1 example.test # comment and spaces, but not controlled
127.0.0.1	localias.test	#{"controller":"localias"}
`)
	lines := Parse(strings.NewReader(contents))
	require.NotNil(t, lines)
	require.Equal(t, contents, asFile(lines))
}

func TestMeta(t *testing.T) {
	contents := etcfile(`127.0.0.1 localhost #{"controller":"localias", "garbage": "hashtag#"}`)
	lines := Parse(strings.NewReader(contents))
	require.NotNil(t, lines)
	require.Len(t, lines, 1)
	require.NotNil(t, lines[0].Entry)
	require.NotNil(t, lines[0].Entry.Meta)
	require.Equal(t, "localias", lines[0].Entry.Meta.Controller)
}

func etcfile(lines ...string) string {
	for i, l := range lines {
		lines[i] = strings.TrimSpace(l)
	}
	return strings.TrimSpace(strings.Join(lines, "\n")) + "\n"
}
