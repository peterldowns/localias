package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/fatih/color"
	"github.com/peterldowns/testy/assert"

	"github.com/peterldowns/localias/cmd/localias/root"
)

// Exercise main's real error handling and exit codes without forking a daemon.
func TestCLIProcess(t *testing.T) { //nolint:paralleltest // runs the CLI and exits the subprocess
	if os.Getenv("LOCALIAS_TEST_CLI") != "1" {
		t.Skip("subprocess helper")
	}
	separator := slices.Index(os.Args, "--")
	if separator == -1 {
		t.Fatal("missing CLI arguments")
	}
	// Force color even with piped output to catch ANSI escapes in JSON/errors.
	color.NoColor = false
	root.Command.SetArgs(os.Args[separator+1:])
	main()
	os.Exit(0)
}

func runCLI(t *testing.T, stateDir string, args ...string) (stdout, stderr string, runErr error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 20*time.Second)
	defer cancel()
	commandArgs := append([]string{"-test.run=^TestCLIProcess$", "--"}, args...)
	cmd := exec.CommandContext(ctx, os.Args[0], commandArgs...)
	cmd.Dir = stateDir
	cmd.Env = append(os.Environ(),
		"LOCALIAS_TEST_CLI=1",
		"LOCALIAS_CONFIGFILE="+filepath.Join(stateDir, "config.yaml"),
		"XDG_STATE_HOME="+stateDir,
		"NO_COLOR=",
	)
	var out, diagnostics bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &diagnostics
	runErr = cmd.Run()
	return out.String(), diagnostics.String(), runErr
}

func assertJSONDocument(t *testing.T, stdout, expected string) {
	t.Helper()
	decoder := json.NewDecoder(strings.NewReader(stdout))
	decoder.UseNumber()
	var actual, want any
	assert.NoError(t, decoder.Decode(&actual))
	assert.NoError(t, json.NewDecoder(strings.NewReader(expected)).Decode(&want))
	// Compare JSON text after normalization, preserving number and boolean types.
	actualJSON, err := json.Marshal(actual)
	assert.NoError(t, err)
	expectedJSON, err := json.Marshal(want)
	assert.NoError(t, err)
	assert.Equal(t, string(expectedJSON), string(actualJSON))
	assert.Equal(t, true, errors.Is(decoder.Decode(&actual), io.EOF))
	assert.Equal(t, false, strings.Contains(stdout, "\x1b"))
}

func TestListStructuredOutput(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name     string
		contents string
		args     []string
		want     string
	}{
		{
			name: "json flag",
			contents: "embed.movie-tickets.localhost: 14204\n" +
				"api-test.movie-tickets.localhost: 14201\n",
			args: []string{"--json"},
			want: `{"aliases":{"api-test.movie-tickets.localhost":14201,"embed.movie-tickets.localhost":14204}}`,
		},
		{
			name:     "json format",
			contents: "api.test: 14201\n",
			args:     []string{"--format", "json"},
			want:     `{"aliases":{"api.test":14201}}`,
		},
		{
			name: "empty json",
			args: []string{"--json"},
			want: `{"aliases":{}}`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			assert.NoError(t, os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(tc.contents), 0o600))
			stdout, stderr, err := runCLI(t, dir, append([]string{"list"}, tc.args...)...)
			assert.NoError(t, err)
			assert.Equal(t, "", stderr)
			assertJSONDocument(t, stdout, tc.want)
		})
	}
}

func TestListReadsSelectedSavedConfiguration(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	selected := filepath.Join(dir, "selected.yaml")
	assert.NoError(t, os.WriteFile(filepath.Join(dir, "config.yaml"), []byte("api.test: 1000\n"), 0o600))
	assert.NoError(t, os.WriteFile(selected, []byte("api.test: 2000\n"), 0o600))
	for _, port := range []int{2000, 3000} {
		assert.NoError(t, os.WriteFile(selected, fmt.Appendf(nil, "api.test: %d\n", port), 0o600))
		stdout, stderr, err := runCLI(t, dir, "--configfile", selected, "list", "--json")
		assert.NoError(t, err)
		assert.Equal(t, "", stderr)
		assertJSONDocument(t, stdout, fmt.Sprintf(`{"aliases":{"api.test":%d}}`, port))
	}
}

func TestListTextPreservesSavedOrder(t *testing.T) {
	t.Parallel()
	for _, args := range [][]string{{"list"}, {"list", "--format=text"}} {
		dir := t.TempDir()
		assert.NoError(t, os.WriteFile(filepath.Join(dir, "config.yaml"), []byte("z.test: 9002\na.test: 9001\n"), 0o600))
		stdout, stderr, err := runCLI(t, dir, args...)
		assert.NoError(t, err)
		assert.Equal(t, "", stderr)
		assert.Equal(t, true, strings.Contains(stdout, "z.test"))
		assert.Equal(t, true, strings.Contains(stdout, "a.test"))
		assert.Equal(t, true, strings.Index(stdout, "z.test") < strings.Index(stdout, "a.test"))
		assert.Equal(t, true, strings.Contains(stdout, "9001"))
		assert.Equal(t, true, strings.Contains(stdout, "9002"))
		assert.Equal(t, true, strings.Contains(stdout, "\x1b"))
	}
}

func TestStatusOutput(t *testing.T) {
	t.Parallel()
	// A reaped child gives us a known dead PID without signaling a live daemon.
	child := exec.CommandContext(t.Context(), "true")
	assert.NoError(t, child.Run())
	currentPID := os.Getpid()
	stalePID := child.Process.Pid
	for _, tc := range []struct {
		name  string
		pid   *int
		empty bool
		json  bool
		want  string
	}{
		{name: "stopped", json: true, want: `{"running":false,"pid":null}`},
		{name: "empty pid", empty: true, json: true, want: `{"running":false,"pid":null}`},
		{name: "running", pid: &currentPID, json: true, want: fmt.Sprintf(`{"running":true,"pid":%d}`, currentPID)},
		{name: "stale pid", pid: &stalePID, json: true, want: `{"running":false,"pid":null}`},
		{name: "stopped text", want: "daemon is not running\n"},
		{name: "running text", pid: &currentPID, want: fmt.Sprintf("daemon running with pid %d\n", currentPID)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			if tc.pid != nil || tc.empty {
				pidFile := filepath.Join(dir, "localias", "daemon.pid")
				assert.NoError(t, os.MkdirAll(filepath.Dir(pidFile), 0o700))
				var data []byte
				if tc.pid != nil {
					data = fmt.Appendf(nil, "%d", *tc.pid)
				}
				assert.NoError(t, os.WriteFile(pidFile, data, 0o600))
			}
			args := []string{"status"}
			if tc.json {
				args = append(args, "--json")
			}
			stdout, stderr, err := runCLI(t, dir, args...)
			assert.NoError(t, err)
			assert.Equal(t, "", stderr)
			if tc.json {
				assertJSONDocument(t, stdout, tc.want)
			} else {
				assert.Equal(t, tc.want, stdout)
			}
		})
	}
}

func TestCLIErrorOutput(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		args []string
		want string
	}{
		{name: "unknown format", args: []string{"list", "--format=xml"}, want: "unsupported output format"},
		{name: "unsupported yaml format", args: []string{"list", "--format=yaml"}, want: "unsupported output format"},
		{name: "conflicting flags", args: []string{"list", "--json", "--format=text"}, want: "none of the others can be"},
		{name: "unknown flag before json", args: []string{"list", "--bad", "--json"}, want: "unknown flag"},
		{name: "unknown flag after json", args: []string{"status", "--json", "--bad"}, want: "unknown flag"},
		{name: "extra list argument", args: []string{"list", "--json", "extra"}, want: "unknown command"},
		{name: "extra status argument", args: []string{"status", "--json", "extra"}, want: "unknown command"},
		{name: "invalid saved config", args: []string{"list", "--json"}, want: "failed to load config"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			assert.NoError(t, os.WriteFile(filepath.Join(dir, "config.yaml"), []byte("invalid: [\n"), 0o600))
			stdout, stderr, err := runCLI(t, dir, tc.args...)
			var exitError *exec.ExitError
			assert.Equal(t, true, errors.As(err, &exitError))
			assert.Equal(t, "", stdout)
			assert.Equal(t, true, strings.Contains(stderr, tc.want))
			assert.Equal(t, false, strings.Contains(stderr, "\x1b"))
		})
	}
}

func TestStatusPIDErrors(t *testing.T) {
	t.Parallel()
	for _, pid := range []string{"not-a-pid", "0", "-1", "directory"} {
		t.Run(pid, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			pidFile := filepath.Join(dir, "localias", "daemon.pid")
			assert.NoError(t, os.MkdirAll(filepath.Dir(pidFile), 0o700))
			if pid == "directory" {
				assert.NoError(t, os.Mkdir(pidFile, 0o700))
			} else {
				assert.NoError(t, os.WriteFile(pidFile, []byte(pid), 0o600))
			}
			stdout, stderr, err := runCLI(t, dir, "status", "--json")
			var exitError *exec.ExitError
			assert.Equal(t, true, errors.As(err, &exitError))
			assert.Equal(t, "", stdout)
			assert.Equal(t, true, strings.Contains(stderr, "daemon PID"))
			assert.Equal(t, false, strings.Contains(stderr, "\x1b"))
		})
	}
}
