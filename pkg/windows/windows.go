package windows

import (
	"errors"
	"runtime"
)

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

func InstallCert(certPath string) error {
	_, err := powershell("scripts/install-cert.ps1", certPath)
	return err
}

func powershell(scriptPath string, args ...string) (string, error) {
	_, _ = scriptPath, args
	return "", errors.New("not implemented") // TODO(hs): implement this with PowerShell invocation too?
}
