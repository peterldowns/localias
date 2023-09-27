package shared

import (
	"fmt"

	"github.com/peterldowns/localias/pkg/config"
	"github.com/peterldowns/localias/pkg/hostctl"
	"github.com/peterldowns/localias/pkg/windows"
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
	cfg, err := config.Load(Flags.Configfile)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}
	return cfg
}

func Controller() hostctl.Controller {
	name := "localias"

	// On Windows, always edit the Windows hosts file
	if windows.IsWindows() {
		return hostctl.NewWindowsController(name)
	}

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
