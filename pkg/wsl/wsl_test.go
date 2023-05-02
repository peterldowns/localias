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

func TestExample(t *testing.T) {
	msg, err := Example("hello, world")
	require.NoError(t, err)
	require.Equal(t, "Received message=hello, world", msg)
}

func TestReadWindowsHosts(t *testing.T) {
	hosts, err := ReadWindowsHosts()
	require.NoError(t, err)
	require.Equal(t, "hello", hosts)
}
