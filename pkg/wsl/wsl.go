package wsl

import (
	"os"
)

func IsWSL() bool {
	msg, err := bash("scripts/is-wsl.sh")
	if err != nil {
		panic(err)
	}
	return msg != ""
}

func IP() string {
	ip, err := bash("scripts/ip.sh")
	if err != nil {
		panic(err)
	}
	return ip
}

func ReadWindowsHosts() (string, error) {
	return bash("scripts/read-windows-hosts.sh")
}

func WriteWindowsHostsFromFile(tmpFilePath string) error {
	winTmpFilePath, err := execute("wslpath", nil, "-w", tmpFilePath)
	if err != nil {
		return err
	}
	_, err = powershell("scripts/write-file.ps1", winTmpFilePath, `$env:windir\System32\drivers\etc\hosts`, "sudo")
	return err
}

func WriteWindowsHosts(contents string) error {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "localias-write-windows-hosts-*")
	if err != nil {
		return err
	}
	defer tmpfile.Close()
	if _, err := tmpfile.Write([]byte(contents)); err != nil {
		return err
	}
	if err := tmpfile.Close(); err != nil {
		return err
	}
	path := tmpfile.Name()
	defer os.Remove(path) // delete the temporary file after the command is done
	return WriteWindowsHostsFromFile(path)
}

func InstallCert(certPath string) error {
	winCertPath, err := execute("wslpath", nil, "-w", certPath)
	if err != nil {
		return err
	}
	_, err = powershell("scripts/install-cert.ps1", winCertPath)
	return err
}
