package wsl

import (
	"bytes"
	"embed"
	"fmt"
	"os/exec"
)

//go:embed scripts
var scripts embed.FS

func IsWSL() bool {
	return false
}

func Bash(script string, args ...string) ([]byte, error) {
	scriptContents, err := scripts.ReadFile("scripts/" + script)
	if err != nil {
		return nil, err
	}

	// this is how you send flags/params/args to a script being executed via
	// stdin:
	//
	// 		cat script.sh | bash -s - foo bar
	//
	// https://stackoverflow.com/a/8514318
	cmdArgs := append([]string{"-s", "-"}, args...)
	cmd := exec.Command("bash", cmdArgs...)

	var stdout, stderr bytes.Buffer
	cmd.Stdin = bytes.NewReader(scriptContents)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		if errMsg := stderr.String(); errMsg != "" {
			return nil, fmt.Errorf("%w: script %s failed: %s", err, script, errMsg)
		}
		return nil, err
	}
	return stdout.Bytes(), nil
}

/*
	cmd := exec.Command("bash")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stdin.Close()

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
		return
	}

	script, err := helloScript.ReadFile("hello.sh")
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(os.Args) > 1 {
		script = []byte(fmt.Sprintf("%s %s", string(script), os.Args[1]))
	}

	_, err = stdin.Write(script)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = stdin.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Println(err)
		return
	}*/
