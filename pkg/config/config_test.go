package config

import (
	"testing"

	"github.com/peterldowns/testy/assert"
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
	assert.NoError(t, err)
	assert.Equal(t, "./example.roundtrip.yaml", cfg.Path)
	assert.Equal(t, exampleEntries, cfg.Entries)
}

func TestWriteConfig(t *testing.T) { //nolint:paralleltest // weird race on the file
	cfg := &Config{
		Path:    "./example.roundtrip.yaml",
		Entries: exampleEntries,
	}
	err := cfg.Save()
	assert.NoError(t, err)
}

func TestConfigRoundtripsPreservingOrder(t *testing.T) { //nolint:paralleltest // weird race on the file
	cfg, err := Open("./example.roundtrip.yaml")
	assert.NoError(t, err)

	err = cfg.Save()
	assert.NoError(t, err)

	cfg2, err := Open(cfg.Path)
	assert.NoError(t, err)
	assert.Equal(t, cfg.Path, cfg2.Path)
	assert.Equal(t, cfg.Entries, cfg2.Entries)
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
	assert.Equal(t, expected, cfg.Entries)

	assert.NoError(t, cfg.Save())
	cfg2, err := Open(cfg.Path)
	assert.NoError(t, err)
	assert.Equal(t, expected, cfg2.Entries)
}

func TestDefaultPath(t *testing.T) {
	t.Parallel()
	path, err := Path(nil)
	assert.NoError(t, err)
	assert.NotEqual(t, "", path)
}

func TestImport(t *testing.T) {
	t.Parallel()
	cfg := &Config{
		Entries: []Entry{
			{Alias: "a", Port: 1},
			{Alias: "b", Port: 2},
		},
	}
	other := &Config{
		Entries: []Entry{
			{Alias: "b", Port: 3}, // will update the existing entry
			{Alias: "c", Port: 4}, // will be a new addition
		},
	}
	added, updated := cfg.Import(other)
	assert.Equal(t, []Entry{{Alias: "c", Port: 4}}, added)
	assert.Equal(t, []Entry{{Alias: "b", Port: 3}}, updated)
	expected := []Entry{
		{Alias: "a", Port: 1},
		{Alias: "b", Port: 3},
		{Alias: "c", Port: 4},
	}
	assert.Equal(t, expected, cfg.Entries)
}
