package hostctl

import (
	"strings"
	"testing"

	"github.com/peterldowns/testy/assert"
)

func TestParseUncontrolledRoundtrips(t *testing.T) {
	t.Parallel()
	contents := etcfile(`
# just a comment     


127.0.0.1	localhost	host.test
`)
	lines := Parse(strings.NewReader(contents))
	assert.NotNil(t, lines)
	assert.Equal(t, 4, len(lines))
	assert.Equal(t, contents, asFile(lines))
}

func TestParseControlled(t *testing.T) {
	t.Parallel()
	contents := etcfile(`
# /etc/hosts file blah blah blah
127.0.0.1 localhost


127.0.0.1 example.test # comment and spaces, but not controlled
127.0.0.1	localias.test	#{"controller":"localias"}
`)
	lines := Parse(strings.NewReader(contents))
	assert.NotNil(t, lines)
	assert.Equal(t, contents, asFile(lines))
}

func TestMeta(t *testing.T) {
	t.Parallel()
	contents := etcfile(`127.0.0.1 localhost #{"controller":"localias", "garbage": "hashtag#"}`)
	lines := Parse(strings.NewReader(contents))
	assert.NotNil(t, lines)
	assert.Equal(t, 1, len(lines))
	assert.NotNil(t, lines[0].Entry)
	assert.NotNil(t, lines[0].Entry.Meta)
	assert.Equal(t, "localias", lines[0].Entry.Meta.Controller)
}

func etcfile(lines ...string) string {
	for i, l := range lines {
		lines[i] = strings.TrimSpace(l)
	}
	return strings.TrimSpace(strings.Join(lines, "\n")) + "\n"
}
