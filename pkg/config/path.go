package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Integralist/go-findroot/find"
	"github.com/adrg/xdg"
	"go.yaml.in/yaml/v2"
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
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, configFileMode)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	contents, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	// A map slice keeps file order so replacements stay in their original position.
	var entries yaml.MapSlice
	if err := yaml.Unmarshal(contents, &entries); err != nil {
		return nil, err
	}
	c := Config{Path: path}
	// YAML values may have types other than the supported integer/string targets.
	// Check keys and parse targets explicitly so errors identify the offending
	// alias instead of panicking on a type assertion.
	for _, entry := range entries {
		alias, ok := entry.Key.(string)
		if !ok || alias == "" {
			return nil, fmt.Errorf("invalid alias %#v: expected a nonempty string", entry.Key)
		}
		target, err := ParseTarget(entry.Value)
		if err != nil {
			return nil, fmt.Errorf("alias %q: invalid target %v: %w", alias, entry.Value, err)
		}
		target.Alias = alias
		c.Set(target)
	}
	return &c, nil
}
