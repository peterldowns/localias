package server

import (
	"testing"

	"github.com/peterldowns/testy/check"
)

func TestEnsureSuffix(t *testing.T) {
	t.Parallel()
	result := ensureSuffix("hostname.local", ".local")
	check.Equal(t, "hostname.local", result)

	result = ensureSuffix("hostname", ".local")
	check.Equal(t, "hostname.local", result)

	result = ensureSuffix("hostname.local.foo.", ".local")
	check.Equal(t, "hostname.local.foo..local", result)
}
