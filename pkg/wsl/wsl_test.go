package wsl

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDetect(t *testing.T) {
	stdout, err := Bash("detect.sh")
	require.NoError(t, err)
	require.NotEqual(t, "", string(stdout))
}

func TestIP(t *testing.T) {
	ip, err := Bash("wsl-ip.sh")
	require.NoError(t, err)
	require.Equal(t, "172.20.166.118", strings.TrimSpace(string(ip)))
}
