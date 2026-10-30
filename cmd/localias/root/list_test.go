package root

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/peterldowns/testy/assert"

	"github.com/peterldowns/localias/pkg/config"
)

func TestJSONOutputErrors(t *testing.T) {
	t.Parallel()
	err := outputJSON(failingWriter{}, []config.Entry{{Alias: "api.test", Port: 14201}})
	assert.Equal(t, true, errors.Is(err, io.ErrClosedPipe))
}

func TestJSONEscapesAliases(t *testing.T) {
	t.Parallel()
	var out bytes.Buffer
	err := outputJSON(&out, []config.Entry{{Alias: "quoted\".test", Port: 14201}})
	assert.NoError(t, err)
	assert.Equal(t, "{\"aliases\":{\"quoted\\\".test\":14201}}\n", out.String())
}

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) {
	return 0, io.ErrClosedPipe
}
