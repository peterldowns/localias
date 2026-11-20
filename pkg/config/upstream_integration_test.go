package config

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddypki"
)

const upstreamRequestTimeout = 5 * time.Second

// TestExplicitUpstreamProxy exercises the generated configuration in the same
// embedded Caddy runtime used by the server, including the upstream Host header.
func TestExplicitUpstreamProxy(t *testing.T) { //nolint:paralleltest // Caddy and XDG state are process globals
	if runIsolated(t) {
		return
	}
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, r.Host)
	}))
	defer upstream.Close()
	checkExplicitUpstreamProxy(t, strings.TrimPrefix(upstream.URL, "http://"))
}

// Run via integration/sibling-container/compose.yaml to exercise Docker DNS.
func TestSiblingContainerUpstreamProxy(t *testing.T) { //nolint:paralleltest // Caddy and XDG state are process globals
	upstream := os.Getenv("LOCALIAS_TEST_SIBLING_UPSTREAM")
	if upstream == "" {
		t.Skip("run the sibling-container integration using integration/sibling-container/compose.yaml")
	}
	if runIsolated(t) {
		return
	}
	checkExplicitUpstreamProxy(t, upstream)
}

// XDG paths are initialized before tests start, so isolate them in a child process.
func runIsolated(t *testing.T) bool {
	t.Helper()
	if os.Getenv("LOCALIAS_TEST_ISOLATED") == "1" {
		return false
	}
	dir := t.TempDir()
	cmd := exec.CommandContext(t.Context(), os.Args[0], "-test.run=^"+t.Name()+"$", "-test.v")
	// The marker lets the child execute the test body without launching another
	// subprocess. Set all XDG directories before package initialization runs.
	cmd.Env = append(os.Environ(), "LOCALIAS_TEST_ISOLATED=1", "XDG_STATE_HOME="+dir, "XDG_DATA_HOME="+dir, "XDG_CONFIG_HOME="+dir, "XDG_CACHE_HOME="+dir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("isolated test: %v\n%s", err, output)
	}
	t.Logf("isolated test passed:\n%s", output)
	return true
}

func checkExplicitUpstreamProxy(t *testing.T, upstream string) {
	t.Helper()

	// Let the OS select a high port, then let Caddy bind it.
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	if err := listener.Close(); err != nil {
		t.Fatal(err)
	}
	host := fmt.Sprintf("proxy.test:%d", port)
	cfg := Config{Entries: []Entry{{Alias: "http://" + host, Upstream: upstream}}}
	data, _, err := cfg.CaddyJSON()
	if err != nil {
		t.Fatalf("generate Caddy JSON: %v", err)
	}
	// The generated configuration enables local PKI. Disable trust installation
	// in this HTTP-only test to avoid modifying the machine's certificate stores.
	var generated caddy.Config
	if err := json.Unmarshal(data, &generated); err != nil {
		t.Fatal(err)
	}
	var pki caddypki.PKI
	if err := json.Unmarshal(generated.AppsRaw["pki"], &pki); err != nil {
		t.Fatalf("decode generated PKI configuration: %v", err)
	}
	ca, ok := pki.CAs["local"]
	if !ok || ca == nil {
		t.Fatal("generated PKI configuration has no local CA")
	}
	ca.InstallTrust = new(false)
	generated.AppsRaw["pki"], err = json.Marshal(pki)
	if err != nil {
		t.Fatal(err)
	}
	data, err = json.Marshal(generated)
	if err != nil {
		t.Fatal(err)
	}
	if err := caddy.Load(data, false); err != nil {
		t.Fatalf("load embedded Caddy: %v", err)
	}
	t.Cleanup(func() {
		if err := caddy.Stop(); err != nil {
			t.Errorf("stop embedded Caddy: %v", err)
		}
	})

	client := &http.Client{Timeout: upstreamRequestTimeout}
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, fmt.Sprintf("http://127.0.0.1:%d/", port), nil)
	if err != nil {
		t.Fatal(err)
	}
	// Dial loopback directly to avoid hosts-file changes, while still selecting
	// the generated alias route and testing the Host sent to the upstream.
	req.Host = host
	response, err := client.Do(req)
	if err != nil {
		t.Fatalf("request through Caddy: %v", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusOK || string(body) != host {
		t.Fatalf("upstream response: status %d, Host %q; want status 200, Host %q", response.StatusCode, body, host)
	}
}
