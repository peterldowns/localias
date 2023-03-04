package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var exampleDirectives = []Directive{ //nolint:gochecknoglobals
	{Alias: "https://secure.local", Port: 9000},
	{Alias: "http://insecure.local", Port: 9001},
	{Alias: "bare.local", Port: 9002},
	{Alias: "bare", Port: 9003},
	{Alias: "invalid://failure", Port: 9004},
	{Alias: "valid.duplicate", Port: 9000},
}

func TestReadConfig(t *testing.T) {
	cfg, err := Open("./example.yaml")
	require.NoError(t, err)
	require.Equal(t, "./example.yaml", cfg.Path)
	require.ElementsMatch(t, exampleDirectives, cfg.Directives)
}

func TestWriteConfig(t *testing.T) {
	cfg := &Config{
		Path:       "./example.yaml",
		Directives: exampleDirectives,
	}
	err := cfg.Save()
	require.NoError(t, err)
}

func TestConfigRoundtrips(t *testing.T) {
	cfg, err := Open("./example.yaml")
	require.NoError(t, err)

	err = cfg.Save()
	require.NoError(t, err)

	cfg2, err := Open(cfg.Path)
	require.NoError(t, err)
	require.Equal(t, cfg.Path, cfg2.Path)
	require.ElementsMatch(t, cfg.Directives, cfg2.Directives)
}

func TestLoad(t *testing.T) {
	cfg, err := Load(nil)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	err = cfg.Save()
	require.NoError(t, err)
}

func TestDefaultPath(t *testing.T) {
	xdgPath, err := DefaultPath()
	require.NoError(t, err)
	require.NotEqual(t, "", xdgPath)
}
