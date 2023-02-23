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
	f := Parse(strings.NewReader(contents))
	require.NotNil(t, f)
	require.Len(t, f.Lines, 4)
	require.Equal(t, contents, f.Contents())
}

func TestParseControlled(t *testing.T) {
	contents := etcfile(`
# /etc/hosts file blah blah blah
127.0.0.1 localhost


127.0.0.1 example.test # comment and spaces, but not controlled
127.0.0.1	pfpro.test	#{"controller":"pfpro"}
`)
	f := Parse(strings.NewReader(contents))
	require.NotNil(t, f)
	require.Equal(t, contents, f.Contents())
}

func TestMeta(t *testing.T) {
	contents := etcfile(`127.0.0.1 localhost #{"controller":"pfpro", "garbage": "hashtag#"}`)
	f := Parse(strings.NewReader(contents))
	require.NotNil(t, f)
	require.Len(t, f.Lines, 1)
	require.NotNil(t, f.Lines[0].Entry)
	require.NotNil(t, f.Lines[0].Entry.Meta)
	require.Equal(t, "pfpro", f.Lines[0].Entry.Meta.Controller)
}

func etcfile(lines ...string) string {
	for i, l := range lines {
		lines[i] = strings.TrimSpace(l)
	}
	return strings.TrimSpace(strings.Join(lines, "\n")) + "\n"
}
