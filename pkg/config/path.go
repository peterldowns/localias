package config

import (
	"io"
	"os"
	"path/filepath"

	"github.com/Integralist/go-findroot/find"
	"github.com/adrg/xdg"
	"github.com/go-yaml/yaml"
)

// Path returns a path to a config file by checking the following list of options.
// The first path that points to an existing file is used.
//
//   - localias --config <path> ...
//   - LOCALIAS_CONFIGFILE=<path> localias ...
//   - .localias.yaml in current directory
//   - .localias.yaml in repository root and command is run from inside a repository
//   - $XDG_CONFIG_HOME/localias.yaml (or OS fallback if XDG_CONFIG_HOME is not set)
func Path(cfgPath *string) (string, error) {
	if cfgPath != nil && *cfgPath != "" {
		return *cfgPath, nil
	}
	if path := os.Getenv("LOCALIAS_CONFIGFILE"); path != "" {
		return path, nil
	}
	if path := lookup("./.localias.yaml"); path != "" {
		return path, nil
	}
	if repo, err := find.Repo(); err == nil {
		if path := lookup(filepath.Join(repo.Path, ".localias.yaml")); path != "" {
			return path, nil
		}
	}
	return xdg.ConfigFile("localias.yaml")
}

// lookup returns the given path if the file exists, otherwise returns empty
// string.
func lookup(path string) string {
	if _, err := os.Stat(path); err != nil {
		return ""
	}
	return path
}

// Open reads and parses a [Config] struct from a yaml file on disk.
func Open(path string) (*Config, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	contents, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	// Use a MapSlice in order to preserve order
	var entries yaml.MapSlice
	if err := yaml.Unmarshal(contents, &entries); err != nil {
		return nil, err
	}
	c := Config{Path: path}
	for _, entry := range entries {
		c.Set(Entry{
			Alias: entry.Key.(string),
			Port:  entry.Value.(int),
		})
	}
	return &c, nil
}
