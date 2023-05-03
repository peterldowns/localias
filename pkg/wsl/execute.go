package wsl

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/alessio/shellescape"
)

//go:embed scripts
var scripts embed.FS

func powershell(scriptPath string, args ...string) (string, error) { //nolint:unparam // stdout not yet used
	scriptContents, err := scripts.ReadFile(scriptPath)
	if err != nil {
		return "", err
	}
	// when invoking a powershell script over stdin, the entire script and its
	// arguments must be wrapped by a script block
	//
	//		& {
	//		<script>
	//		} 'arg1' ... 'argN'
	//
	// additionally, the entire package must be delivered with an extra newline
	// at the end or the execution will be ignored and the invocation will exit
	// with status code 0 but not run your script at all (!)
	//
	// see https://stackoverflow.com/a/42475326
	wrappedScript := fmt.Sprintf(
		"& {\n%s} %s\n\n",
		scriptContents,
		shellescape.QuoteCommand(args),
	)
	return execute("powershell.exe", strings.NewReader(wrappedScript), "-command", "-")
}

func bash(scriptPath string, args ...string) (string, error) {
	scriptContents, err := scripts.ReadFile(scriptPath)
	if err != nil {
		return "", err
	}

	// this is how you send flags/params/args to a script being executed via
	// stdin:
	//
	// 		cat script.sh | bash -s - foo bar
	//
	// https://stackoverflow.com/a/8514318
	cmdArgs := append([]string{"-s", "-"}, args...)
	return execute("bash", bytes.NewReader(scriptContents), cmdArgs...)
}

func execute(program string, stdin io.Reader, args ...string) (string, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(program, args...)
	cmd.Stdin = stdin
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return "", err
	}

	if err := cmd.Wait(); err != nil {
		if errMsg := stderr.String(); errMsg != "" {
			return "", fmt.Errorf("program %s failed with error(%w): %s", program, err, errMsg)
		}
		return "", fmt.Errorf("program %s failed with error(%w)", program, err)
	}
	return strings.TrimSuffix(stdout.String(), "\n"), nil
}
