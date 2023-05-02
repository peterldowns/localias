package wsl

func IsWSL() bool {
	msg, err := bash("scripts/detect.sh")
	if err != nil {
		panic(err)
	}
	return msg != ""
}

func IP() string {
	ip, err := bash("scripts/wsl-ip.sh")
	if err != nil {
		panic(err)
	}
	return ip
}

func ReadWindowsHosts() (string, error) {
	return bash("scripts/read-etc-hosts.sh")
}

func Example(message string) (string, error) {
	out, err := powershell("scripts/example.ps1", message)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
