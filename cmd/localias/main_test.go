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
			name:     "mixed targets",
			contents: "web.test: 5173\napi.test: backend:8080\n",
			args:     []string{"--json"},
			want:     `{"aliases":{"web.test":5173,"api.test":"backend:8080"}}`,
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
			_, isExitError := errors.AsType[*exec.ExitError](err)
			assert.Equal(t, true, isExitError)
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
			_, isExitError := errors.AsType[*exec.ExitError](err)
			assert.Equal(t, true, isExitError)
			assert.Equal(t, "", stdout)
			assert.Equal(t, true, strings.Contains(stderr, "daemon PID"))
			assert.Equal(t, false, strings.Contains(stderr, "\x1b"))
		})
	}
}

func TestSetTargets(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	for _, tc := range []struct {
		args []string
		want string
	}{
		{[]string{"set", "api.test", "9000"}, `{"aliases":{"api.test":9000}}`},
		{[]string{"set", "api.test", "backend:8080"}, `{"aliases":{"api.test":"backend:8080"}}`},
		{[]string{"set", "--alias", "api.test", "--port", "9001"}, `{"aliases":{"api.test":9001}}`},
		{[]string{"set", "--alias", "api.test", "--upstream", "[::1]:8080"}, `{"aliases":{"api.test":"[::1]:8080"}}`},
	} {
		stdout, stderr, err := runCLI(t, dir, tc.args...)
		assert.NoError(t, err)
		assert.Equal(t, "", stderr)
		assert.Equal(t, true, strings.Contains(stdout, tc.args[len(tc.args)-1]))
		stdout, stderr, err = runCLI(t, dir, "list", "--json")
		assert.NoError(t, err)
		assert.Equal(t, "", stderr)
		assertJSONDocument(t, stdout, tc.want)
	}
	stdout, _, err := runCLI(t, dir, "list")
	assert.NoError(t, err)
	assert.Equal(t, true, strings.Contains(stdout, "[::1]:8080"))
	stdout, _, err = runCLI(t, dir, "rm", "api.test")
	assert.NoError(t, err)
	assert.Equal(t, true, strings.Contains(stdout, "[::1]:8080"))
}

func TestSetRejectsInvalidArguments(t *testing.T) {
	t.Parallel()
	for _, args := range [][]string{
		{"set"},
		{"set", "api.test"},
		{"set", "api.test", "9000", "extra"},
		{"set", "api.test", "9000", "--alias", ""},
		{"set", "api.test", "9000", "--port", "0"},
		{"set", "--alias", "api.test"},
		{"set", "--port", "9000"},
		{"set", "--upstream", "backend:8080"},
		{"set", "api.test", "9000", "--upstream", "backend:8080"},
		{"set", "--alias", "api.test", "--port", "9000", "--upstream", "backend:8080"},
		{"set", "--alias", "api.test", "--port", "0"},
		{"set", "--alias", "api.test", "--upstream", ""},
		{"set", "--alias", "api.test", "--upstream", "9000"},
		{"set", "--alias", "api.test", "--upstream", "00080"},
		{"set", "api.test", "https://backend:8080"},
	} {
		t.Run(strings.Join(args, "_"), func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			path := filepath.Join(dir, "config.yaml")
			original := []byte("existing.test: 1234\n")
			assert.NoError(t, os.WriteFile(path, original, 0o600))
			stdout, stderr, err := runCLI(t, dir, args...)
			_, isExitError := errors.AsType[*exec.ExitError](err)
			assert.Equal(t, true, isExitError)
			assert.Equal(t, "", stdout)
			assert.Equal(t, true, stderr != "" && !strings.Contains(stderr, "panic:"))
			contents, readErr := os.ReadFile(path)
			assert.NoError(t, readErr)
			assert.Equal(t, string(original), string(contents))
		})
	}
}

func TestImportAndClearUpstream(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "import.yaml")
	assert.NoError(t, os.WriteFile(path, []byte("api.test: backend:8080\n"), 0o600))
	stdout, stderr, err := runCLI(t, dir, "import", path)
	assert.NoError(t, err)
	assert.Equal(t, "", stderr)
	assert.Equal(t, true, strings.Contains(stdout, "backend:8080"))
	stdout, stderr, err = runCLI(t, dir, "list", "--json")
	assert.NoError(t, err)
	assert.Equal(t, "", stderr)
	assertJSONDocument(t, stdout, `{"aliases":{"api.test":"backend:8080"}}`)
	stdout, _, err = runCLI(t, dir, "clear")
	assert.NoError(t, err)
	assert.Equal(t, true, strings.Contains(stdout, "backend:8080"))
}

func TestQuotedNumericTargetNormalizesOnSave(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	assert.NoError(t, os.WriteFile(path, []byte("api.test: \"9000\"\n"), 0o600))
	stdout, _, err := runCLI(t, dir, "list", "--json")
	assert.NoError(t, err)
	assertJSONDocument(t, stdout, `{"aliases":{"api.test":9000}}`)
	_, _, err = runCLI(t, dir, "set", "other.test", "9001")
	assert.NoError(t, err)
	contents, err := os.ReadFile(path)
	assert.NoError(t, err)
	assert.Equal(t, true, strings.Contains(string(contents), "api.test: 9000\n"))
	stdout, _, err = runCLI(t, dir, "list", "--json")
	assert.NoError(t, err)
	assertJSONDocument(t, stdout, `{"aliases":{"api.test":9000,"other.test":9001}}`)
}

func TestMutatingCommandsRejectInvalidSavedConfig(t *testing.T) {
	t.Parallel()
	for _, contents := range []string{"invalid: [\n", "api.test: https://backend:8080\n"} {
		for _, args := range [][]string{{"import", "import.yaml"}, {"remove", "api.test"}, {"clear"}} {
			t.Run(args[0]+"/"+contents, func(t *testing.T) {
				t.Parallel()
				dir := t.TempDir()
				path := filepath.Join(dir, "config.yaml")
				assert.NoError(t, os.WriteFile(path, []byte(contents), 0o600))
				assert.NoError(t, os.WriteFile(filepath.Join(dir, "import.yaml"), []byte("other.test: 9000\n"), 0o600))
				stdout, stderr, err := runCLI(t, dir, args...)
				_, isExitError := errors.AsType[*exec.ExitError](err)
				assert.Equal(t, true, isExitError)
				assert.Equal(t, "", stdout)
				assert.Equal(t, true, strings.Contains(stderr, "failed to load config"))
				assert.Equal(t, false, strings.Contains(stderr, "panic:"))
				after, readErr := os.ReadFile(path)
				assert.NoError(t, readErr)
				assert.Equal(t, contents, string(after))
			})
		}
	}
}

func TestSetUpstreamErrorMessages(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		args []string
		want string
	}{
		{
			name: "empty upstream",
			args: []string{"set", "--alias", "api.test", "--upstream", ""},
			want: "invalid target: target must not be empty",
		},
		{
			name: "numeric upstream",
			args: []string{"set", "--alias", "api.test", "--upstream", "9000"},
			want: "invalid target: upstream must be a host:port address",
		},
		{
			name: "leading zero numeric upstream",
			args: []string{"set", "--alias", "api.test", "--upstream", "00080"},
			want: "invalid target: upstream must be a host:port address",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			stdout, stderr, err := runCLI(t, dir, tc.args...)
			_, isExitError := errors.AsType[*exec.ExitError](err)
			assert.Equal(t, true, isExitError)
			assert.Equal(t, "", stdout)
			assert.Equal(t, true, strings.Contains(stderr, tc.want))
			assert.Equal(t, false, strings.Contains(stderr, "panic:"))
		})
	}
}

func TestImportValidatesAllFilesBeforeSavingOrPrinting(t *testing.T) {
	t.Parallel()
	for _, invalidContents := range []string{"invalid: [\n", "bad.test: https://backend:8080\n"} {
		t.Run(invalidContents, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			configPath := filepath.Join(dir, "config.yaml")
			original := "existing.test: 1234\n"
			assert.NoError(t, os.WriteFile(configPath, []byte(original), 0o600))
			assert.NoError(t, os.WriteFile(filepath.Join(dir, "valid.yaml"), []byte("valid.test: 9000\n"), 0o600))
			assert.NoError(t, os.WriteFile(filepath.Join(dir, "invalid.yaml"), []byte(invalidContents), 0o600))

			stdout, stderr, err := runCLI(t, dir, "import", "valid.yaml", "invalid.yaml")
			_, isExitError := errors.AsType[*exec.ExitError](err)
			assert.Equal(t, true, isExitError)
			assert.Equal(t, "", stdout)
			assert.Equal(t, true, stderr != "" && !strings.Contains(stderr, "panic:"))

			after, readErr := os.ReadFile(configPath)
			assert.NoError(t, readErr)
			assert.Equal(t, original, string(after))
		})
	}
}

func TestImportMultipleFilesPreservesOutputOrder(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	assert.NoError(t, os.WriteFile(filepath.Join(dir, "config.yaml"), []byte("a.test: 9000\n"), 0o600))
	assert.NoError(t, os.WriteFile(filepath.Join(dir, "first.yaml"), []byte("b.test: 9001\na.test: 9002\n"), 0o600))
	assert.NoError(t, os.WriteFile(filepath.Join(dir, "second.yaml"), []byte("c.test: backend:8080\nb.test: 9003\n"), 0o600))

	stdout, stderr, err := runCLI(t, dir, "import", "first.yaml", "second.yaml")
	assert.NoError(t, err)
	assert.Equal(t, "", stderr)
	lines := strings.Split(strings.TrimSuffix(stdout, "\n"), "\n")
	assert.Equal(t, 4, len(lines))
	expected := []struct {
		action string
		alias  string
		target string
	}{
		{"[added]", "b.test", "9001"},
		{"[updated]", "a.test", "9002"},
		{"[added]", "c.test", "backend:8080"},
		{"[updated]", "b.test", "9003"},
	}
	for i, want := range expected {
		assert.Equal(t, true, strings.Contains(lines[i], want.action))
		assert.Equal(t, true, strings.Contains(lines[i], want.alias))
		assert.Equal(t, true, strings.Contains(lines[i], want.target))
	}

	listOut, listErrOut, listErr := runCLI(t, dir, "list", "--json")
	assert.NoError(t, listErr)
	assert.Equal(t, "", listErrOut)
	assertJSONDocument(t, listOut, `{"aliases":{"a.test":9002,"b.test":9003,"c.test":"backend:8080"}}`)
}

func TestRuntimeAndDebugCommandsRejectInvalidSavedConfig(t *testing.T) {
	t.Parallel()
	for _, contents := range []string{"invalid: [\n", "api.test: https://backend:8080\n"} {
		for _, tc := range []struct {
			name    string
			args    []string
			livePID bool
		}{
			{name: "debug_caddyfile", args: []string{"debug", "caddyfile"}},
			{name: "debug_hostctl_apply", args: []string{"debug", "hostctl", "apply"}},
			{name: "run", args: []string{"run"}},
			{name: "start", args: []string{"start"}},
			{name: "start_with_live_pid", args: []string{"start"}, livePID: true},
		} {
			t.Run(tc.name+"/"+contents, func(t *testing.T) {
				t.Parallel()
				dir := t.TempDir()
				path := filepath.Join(dir, "config.yaml")
				assert.NoError(t, os.WriteFile(path, []byte(contents), 0o600))
				if tc.livePID {
					pidFile := filepath.Join(dir, "localias", "daemon.pid")
					assert.NoError(t, os.MkdirAll(filepath.Dir(pidFile), 0o700))
					assert.NoError(t, os.WriteFile(pidFile, fmt.Appendf(nil, "%d", os.Getpid()), 0o600))
				}

				stdout, stderr, err := runCLI(t, dir, tc.args...)
				_, isExitError := errors.AsType[*exec.ExitError](err)
				assert.Equal(t, true, isExitError)
				assert.Equal(t, "", stdout)
				assert.Equal(t, true, strings.Contains(stderr, "failed to load config"))
				assert.Equal(t, false, strings.Contains(stderr, "panic:"))
			})
		}
	}
}

func TestDebugConfigAllowsInvalidAndMissingConfig(t *testing.T) {
	t.Parallel()
	for _, contents := range []string{"invalid: [\n", "api.test: https://backend:8080\n"} {
		t.Run(contents, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			path := filepath.Join(dir, "config.yaml")
			assert.NoError(t, os.WriteFile(path, []byte(contents), 0o600))

			stdout, stderr, err := runCLI(t, dir, "debug", "config")
			assert.NoError(t, err)
			assert.Equal(t, "", stderr)
			assert.Equal(t, path+"\n", stdout)

			stdout, stderr, err = runCLI(t, dir, "debug", "config", "--print")
			assert.NoError(t, err)
			assert.Equal(t, "", stderr)
			assert.Equal(t, contents+"\n", stdout)
		})
	}

	t.Run("missing_config", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		path := filepath.Join(dir, "config.yaml")

		stdout, stderr, err := runCLI(t, dir, "debug", "config")
		assert.NoError(t, err)
		assert.Equal(t, "", stderr)
		assert.Equal(t, path+"\n", stdout)
		_, statErr := os.Stat(path)
		assert.Equal(t, true, errors.Is(statErr, os.ErrNotExist))
	})
}

func TestDebugCertAllowsInvalidAndMissingConfig(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name     string
		contents string
		create   bool
	}{
		{name: "invalid_yaml", contents: "invalid: [\n", create: true},
		{name: "invalid_target", contents: "api.test: https://backend:8080\n", create: true},
		{name: "missing_config", create: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			configPath := filepath.Join(dir, "config.yaml")
			if tc.create {
				assert.NoError(t, os.WriteFile(configPath, []byte(tc.contents), 0o600))
			}

			stdout, stderr, err := runCLI(t, dir, "debug", "cert")
			assert.NoError(t, err)
			assert.Equal(t, "", stderr)
			wantCertPath := filepath.Join(dir, "localias", "caddy", "pki", "authorities", "local", "root.crt")
			assert.Equal(t, wantCertPath+"\n", stdout)
			if !tc.create {
				_, statErr := os.Stat(configPath)
				assert.Equal(t, true, errors.Is(statErr, os.ErrNotExist))
			}
		})
	}
}
