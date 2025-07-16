//go:build manual

// These are tests designed to be run manually while working in a WSL environment.
// They are not included in the automatic test suite since they depend on being run
// from inside a WSL environment.

package wsl

import (
	"testing"

	"github.com/peterldowns/testy/assert"
)

func TestDetect(t *testing.T) {
	assert.Equal(t, true, IsWSL())
}

func TestIP(t *testing.T) {
	ip := IP()
	assert.Equal(t, "172.20.166.118", ip)
}

func TestReadWindowsHosts(t *testing.T) {
	hosts, err := ReadWindowsHosts()
	assert.NoError(t, err)
	assert.NotEqual(t, "", hosts)
}

func TestWriteWindowsHosts(t *testing.T) {
	hosts, err := ReadWindowsHosts()
	assert.NoError(t, err)
	hosts += "\n# added from inside golang TestWriteWindowsHosts!"
	err = WriteWindowsHosts(hosts)
	assert.NoError(t, err)
}

func TestInstallCert(t *testing.T) {
	path := "/home/pd/.local/state/localias/caddy/pki/authorities/local/root.crt"
	err := InstallCert(path)
	assert.NoError(t, err)
}
