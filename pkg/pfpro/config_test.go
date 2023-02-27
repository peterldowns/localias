package pfpro

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var exampleDirectives = []Directive{ //nolint:gochecknoglobals
	{Upstream: "https://secure.local", Downstream: "9000"},
	{Upstream: "http://insecure.local", Downstream: "9001"},
	{Upstream: "bare.local", Downstream: "9002"},
	{Upstream: "bare", Downstream: "9003"},
	{Upstream: "invalid://failure", Downstream: "9004"},
	{Upstream: "valid.duplicate", Downstream: "9000"},
}

func TestReadConfig(t *testing.T) {
	cfg, err := ReadConfig("./example.yaml")
	require.NoError(t, err)
	require.Equal(t, "./example.yaml", cfg.Path)
	require.ElementsMatch(t, exampleDirectives, cfg.Directives)
}

func TestWriteConfig(t *testing.T) {
	cfg := &Config{
		Path:       "./example.yaml",
		Directives: exampleDirectives,
	}
	err := WriteConfig(cfg)
	require.NoError(t, err)
}

func TestConfigRoundtrips(t *testing.T) {
	cfg, err := ReadConfig("./example.yaml")
	require.NoError(t, err)

	err = WriteConfig(cfg)
	require.NoError(t, err)

	cfg2, err := ReadConfig(cfg.Path)
	require.NoError(t, err)
	require.Equal(t, cfg.Path, cfg2.Path)
	require.ElementsMatch(t, cfg.Directives, cfg2.Directives)
}

func TestConfigPath(t *testing.T) {
	xdgPath, err := DefaultConfigPath()
	require.NoError(t, err)
	require.NotEqual(t, "", xdgPath)
}

func TestLoad(t *testing.T) {
	cfg, err := Load(nil)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	err = WriteConfig(cfg)
	require.NoError(t, err)
}
