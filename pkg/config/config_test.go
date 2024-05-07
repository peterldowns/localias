package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var exampleEntries = []Entry{ //nolint:gochecknoglobals
	{Alias: "bare", Port: 9003},
	{Alias: "bare.test", Port: 9002},
	{Alias: "invalid://failure", Port: 9004},
	{Alias: "valid.duplicate", Port: 9000},
	{Alias: "http://insecure.test", Port: 9001},
	{Alias: "https://secure.test", Port: 9000},
}

func TestReadConfig(t *testing.T) {
	t.Parallel()
	cfg, err := Open("./example.roundtrip.yaml")
	require.NoError(t, err)
	require.Equal(t, "./example.roundtrip.yaml", cfg.Path)
	require.ElementsMatch(t, exampleEntries, cfg.Entries)
	require.Equal(t, exampleEntries, cfg.Entries)
}

func TestWriteConfig(t *testing.T) { //nolint:paralleltest // weird race on the file
	cfg := &Config{
		Path:    "./example.roundtrip.yaml",
		Entries: exampleEntries,
	}
	err := cfg.Save()
	require.NoError(t, err)
}

func TestConfigRoundtripsPreservingOrder(t *testing.T) { //nolint:paralleltest // weird race on the file
	cfg, err := Open("./example.roundtrip.yaml")
	require.NoError(t, err)

	err = cfg.Save()
	require.NoError(t, err)

	cfg2, err := Open(cfg.Path)
	require.NoError(t, err)
	require.Equal(t, cfg.Path, cfg2.Path)
	require.Equal(t, cfg.Entries, cfg2.Entries)
}

func TestUpsertUpdatesExistingEntry(t *testing.T) { //nolint:paralleltest // weird race on the file
	cfg := &Config{
		Path: "./example.upsert.yaml",
	}
	cfg.Set(Entry{
		Alias: "dev.test",
		Port:  8000,
	})
	cfg.Set(Entry{
		Alias: "dev.test",
		Port:  9000,
	})
	expected := []Entry{
		{Alias: "dev.test", Port: 9000},
	}
	require.Equal(t, expected, cfg.Entries)

	require.NoError(t, cfg.Save())
	cfg2, err := Open(cfg.Path)
	require.NoError(t, err)
	require.Equal(t, expected, cfg2.Entries)
}

func TestDefaultPath(t *testing.T) {
	t.Parallel()
	path, err := Path(nil)
	require.NoError(t, err)
	require.NotEqual(t, "", path)
}
