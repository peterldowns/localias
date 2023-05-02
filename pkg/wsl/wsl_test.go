//go:build manual

// These are tests designed to be run manually while working in a WSL environment.
// They are not included in the automatic test suite since they depend on being run
// from inside a WSL environment.

package wsl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDetect(t *testing.T) {
	require.Equal(t, true, IsWSL())
}

func TestIP(t *testing.T) {
	ip := IP()
	require.Equal(t, "172.20.166.118", ip)
}

func TestReadWindowsHosts(t *testing.T) {
	hosts, err := ReadWindowsHosts()
	require.NoError(t, err)
	require.NotEqual(t, "", hosts)
}

func TestWriteWindowsHosts(t *testing.T) {
	hosts, err := ReadWindowsHosts()
	require.NoError(t, err)
	hosts += "\n# added from inside golang TestWriteWindowsHosts!"
	err = WriteWindowsHosts(hosts)
	require.NoError(t, err)
}

func TestInstallCert(t *testing.T) {
	path := "/home/pd/.local/state/localias/caddy/pki/authorities/local/root.crt"
	err := InstallCert(path)
	require.NoError(t, err)
}
