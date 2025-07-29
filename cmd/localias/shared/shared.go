package shared

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/peterldowns/localias/pkg/config"
	"github.com/peterldowns/localias/pkg/hostctl"
	"github.com/peterldowns/localias/pkg/wsl"
)

// These will be set at build time with ldflags, see Justfile for how they're
// defined and passed.
var (
	Version = "unknown" //nolint:gochecknoglobals
	Commit  = "unknown" //nolint:gochecknoglobals
)

var Flags struct { //nolint:gochecknoglobals
	Configfile *string
}

func Config() *config.Config {
	path, err := config.Path(Flags.Configfile)
	if err != nil {
		panic(fmt.Errorf("failed to find config: %w", err))
	}
	cfg, err := config.Open(path)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}
	return cfg
}

func Controller() hostctl.Controller {
	name := "localias"
	// On WSL/Mac/Linux, we're always going to need to edit /etc/hosts.
	etcHostController := hostctl.NewFileController(
		"/etc/hosts",
		true,
		name,
	)
	// If we're on WSL, we'll also need to update the window machine's
	// host file.
	if wsl.IsWSL() {
		wslController := hostctl.NewWSLController(name)
		return hostctl.NewMultiController(
			wslController,
			etcHostController,
		)
	}
	return etcHostController
}

func VersionString() string {
	return fmt.Sprintf("%s+commit.%s", Version, Commit)
}

// PrintUpdate is a helper for printing a newly added or updated config entry,
// used by the `set` and `import` commands. It's responsible for pretty colors
// and consistent formatting.
func PrintUpdate(entry config.Entry, updated bool) {
	action := "[added]"
	if updated {
		action = "[updated]"
	}
	fmt.Printf(
		"%s %s -> %s\n",
		color.New(color.FgGreen).Sprint(action),
		color.New(color.FgBlue).Sprintf(entry.Alias),
		color.New(color.FgWhite).Sprintf("%d", entry.Port),
	)
}
