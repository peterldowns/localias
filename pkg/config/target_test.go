package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/adrg/xdg"
)

func TestParseTarget(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input any
		want  Entry
	}{
		{1, Entry{Port: 1}},
		{65535, Entry{Port: 65535}},
		{80, Entry{Port: 80}},
		{"00080", Entry{Port: 80}},
		{"localhost:8080", Entry{Upstream: "localhost:8080"}},
		{"frontend.internal.:5173", Entry{Upstream: "frontend.internal.:5173"}},
		{"service_name.internal:00080", Entry{Upstream: "service_name.internal:80"}},
		{"127.0.0.1:443", Entry{Upstream: "127.0.0.1:443"}},
		{"frontend:65535", Entry{Upstream: "frontend:65535"}},
		{"[::1]:8080", Entry{Upstream: "[::1]:8080"}},
		{"[0:0:0:0:0:0:0:1]:00080", Entry{Upstream: "[::1]:80"}},
		{"[::ABCD]:8080", Entry{Upstream: "[::abcd]:8080"}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprint(tt.input), func(t *testing.T) {
			t.Parallel()
			got, err := ParseTarget(tt.input)
			if err != nil || got != tt.want {
				t.Errorf("ParseTarget(%v) = %#v, %v; want %#v", tt.input, got, err, tt.want)
			}
		})
	}
	invalid := []any{nil, true, 1.5, []any{1}, 0, -1, 65536, "", "0", "65536", "-1", "http://example.com:80", "example.com", ":80", "example.com:http", "::1:80", "[localhost]:80", "[::1%lo0]:80", "foo:80/path", "user@foo:80", "foo bar:80", "foo:80\nrespond hacked", "foo}:80", "foo;bar:80", "foo..bar:80", "999.1.2.3:80", "127.00.0.1:80", "127.0.1:80", "127.1:80", "1.2.3.4.5:80", "12345:80", "service.8080:80", "127.0.0.1.:80", "host:+80", "host:-80", "host: 80", "host:8_0"}
	for _, input := range invalid {
		t.Run("invalid/"+fmt.Sprint(input), func(t *testing.T) {
			t.Parallel()
			if got, err := ParseTarget(input); err == nil {
				t.Errorf("ParseTarget(%v) = %#v; want error", input, got)
			}
		})
	}
}

func TestMixedTargetsRoundTripAndImport(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "config.yaml")
	text := "first: 8080\nsecond: backend.internal:3000\nthird: '9000'\n"
	if err := os.WriteFile(path, []byte(text), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	added, updated := cfg.Import(&Config{Entries: []Entry{{Alias: "first", Upstream: "[::1]:42"}, {Alias: "second", Port: 4321}, {Alias: "fourth", Upstream: "backend.internal:3000"}}})
	if len(added) != 1 || len(updated) != 2 {
		t.Fatalf("added %v updated %v", added, updated)
	}
	want := []Entry{
		{Alias: "first", Upstream: "[::1]:42"},
		{Alias: "second", Port: 4321},
		{Alias: "third", Port: 9000},
		{Alias: "fourth", Upstream: "backend.internal:3000"},
	}
	if !reflect.DeepEqual(cfg.Entries, want) {
		t.Fatalf("imported entries = %#v, want %#v", cfg.Entries, want)
	}
	if err := cfg.Save(); err != nil {
		t.Fatal(err)
	}
	saved, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(saved), "third: 9000") || !strings.Contains(string(saved), "fourth: backend.internal:3000") {
		t.Fatalf("unexpected YAML: %s", saved)
	}
	got, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got.Entries, cfg.Entries) {
		t.Fatalf("round trip = %#v, want %#v", got.Entries, cfg.Entries)
	}
	if got.Entries[0].String() != "first: [::1]:42" || !strings.Contains(got.Entries[0].Caddyfile(), "reverse_proxy [::1]:42") {
		t.Fatal("upstream not displayed or rendered")
	}
}

func TestOpenRejectsInvalidValues(t *testing.T) {
	t.Parallel()
	for _, content := range []string{"42: 8000\n", "\"\": 8000\n", "bad: true\n", "bad: null\n", "bad: [8000]\n", "bad: {port: 8000}\n", "bad: 'host:80/path'\n"} {
		path := filepath.Join(t.TempDir(), "config.yaml")
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
		_, err := Open(path)
		if err == nil || !strings.Contains(err.Error(), "alias") {
			t.Errorf("Open(%q) error = %v", content, err)
			continue
		}
		if content == "\"\": 8000\n" && !strings.Contains(err.Error(), `invalid alias ""`) {
			t.Errorf("Open(%q) error lacks empty alias: %v", content, err)
		}
		if strings.HasPrefix(content, "bad:") && !strings.Contains(err.Error(), `alias "bad": invalid target`) {
			t.Errorf("Open(%q) error lacks alias and target: %v", content, err)
		}
		if strings.Contains(content, "host:80/path") && !strings.Contains(err.Error(), "host:80/path") {
			t.Errorf("error lacks target: %v", err)
		}
	}
}

func TestInvalidProgrammaticTargetsRejected(t *testing.T) {
	t.Parallel()
	for _, entry := range []Entry{{Alias: "bad"}, {Alias: "bad", Port: 80, Upstream: "host:80"}, {Alias: "bad", Port: 65536}, {Alias: "bad", Upstream: "host:80\nrespond hacked"}} {
		cfg := Config{Path: filepath.Join(t.TempDir(), "config.yaml"), Entries: []Entry{entry}}
		saveErr := cfg.Save()
		if saveErr == nil {
			t.Errorf("Save(%#v) succeeded", entry)
		}
		_, _, jsonErr := cfg.CaddyJSON()
		if jsonErr == nil {
			t.Errorf("CaddyJSON(%#v) succeeded", entry)
		}
		for _, err := range []error{saveErr, jsonErr} {
			if err == nil {
				continue
			}
			if entry == (Entry{Alias: "bad"}) && (strings.Contains(err.Error(), `invalid target "0"`) || !strings.Contains(err.Error(), "invalid target:")) {
				t.Errorf("unexpected error for empty target %#v: %v", entry, err)
			}
			if entry == (Entry{Alias: "bad", Port: 80, Upstream: "host:80"}) && (!strings.Contains(err.Error(), "port=80") || !strings.Contains(err.Error(), "host:80")) {
				t.Errorf("error missing both targets for %#v: %v", entry, err)
			}
		}
	}
}

func TestCaddyUpstreamAdaptation(t *testing.T) { //nolint:paralleltest // Caddy tests modify global XDG state
	previousStateHome := xdg.StateHome
	xdg.StateHome = t.TempDir()
	t.Cleanup(func() { xdg.StateHome = previousStateHome })

	cfg := Config{Entries: []Entry{
		{Alias: "port.test", Port: 8000},
		{Alias: "ipv4.test", Upstream: "127.0.0.1:8001"},
		{Alias: "ipv6.test", Upstream: "[::1]:8002"},
		{Alias: "service.test", Upstream: "unresolved_service.internal:8003"},
	}}
	legacy := `port.test {
	reverse_proxy localhost:8000
	tls {
		on_demand
		issuer internal {
			ca local
		}
	}
}`
	if got := cfg.Entries[0].Caddyfile(); got != legacy {
		t.Fatalf("legacy Caddyfile changed:\n%s", got)
	}
	data, _, err := cfg.CaddyJSON()
	if err != nil {
		t.Fatal(err)
	}
	var parsed any
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatal(err)
	}
	var dials []string
	// Caddy groups routes when adapting the file. Find upstream dial addresses
	// throughout the JSON without depending on a particular route nesting order.
	var visit func(any)
	visit = func(v any) {
		switch x := v.(type) {
		case map[string]any:
			for k, value := range x {
				if k == "dial" {
					if dial, ok := value.(string); ok {
						dials = append(dials, dial)
					}
				}
				visit(value)
			}
		case []any:
			for _, value := range x {
				visit(value)
			}
		}
	}
	visit(parsed)
	for _, want := range []string{"localhost:8000", "127.0.0.1:8001", "[::1]:8002", "unresolved_service.internal:8003"} {
		if !slices.Contains(dials, want) {
			t.Errorf("missing dial %q in %v", want, dials)
		}
	}
}

func TestHostnameValidation(t *testing.T) {
	t.Parallel()
	label63 := strings.Repeat("a", 63)
	name253 := strings.Join([]string{label63, label63, label63, strings.Repeat("b", 61)}, ".")
	cases := []struct {
		name, host string
		valid      bool
	}{
		{"single service underscore", "_service", true},
		{"absolute service", "service_name.internal.", true},
		{"maximum label", label63, true},
		{"long label", strings.Repeat("a", 64), false},
		{"maximum name", name253, true},
		{"maximum absolute name", name253 + ".", true},
		{"overlong name", name253 + "b", false},
		{"overlong absolute name", name253 + "b.", false},
		{"root", ".", false},
		{"empty label", "service..internal", false},
		{"leading dot", ".service", false},
		{"double trailing dot", "service..", false},
		{"escaped dot", `service\.internal`, false},
		{"escaped character", `service\065`, false},
		{"caddy brace", "service{internal", false},
		{"caddy quote", `service"internal`, false},
		{"leading hyphen", "-service", false},
		{"trailing hyphen", "service-", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := ParseTarget(tc.host + ":80")
			if (err == nil) != tc.valid {
				t.Errorf("ParseTarget(%q) error = %v; valid = %v", tc.host, err, tc.valid)
			}
		})
	}
}

func TestProgrammaticTargetNormalization(t *testing.T) { //nolint:paralleltest // Caddy tests modify global XDG state
	previousStateHome := xdg.StateHome
	xdg.StateHome = t.TempDir()
	t.Cleanup(func() { xdg.StateHome = previousStateHome })

	raw := Entry{Alias: "frontend.test", Upstream: "frontend:00080"}
	normalized, err := raw.NormalizeTarget()
	if err != nil {
		t.Fatal(err)
	}
	want := Entry{Alias: raw.Alias, Upstream: "frontend:80"}
	if normalized != want {
		t.Fatalf("NormalizeTarget() = %#v, want %#v", normalized, want)
	}
	parsed, err := ParseTarget("frontend:00080")
	if err != nil {
		t.Fatal(err)
	}
	parsed.Alias = raw.Alias
	if parsed != normalized {
		t.Fatalf("parsed = %#v, normalized = %#v", parsed, normalized)
	}
	if raw.Target() != "frontend:80" || raw.UpstreamAddress() != "frontend:80" || !strings.Contains(raw.Caddyfile(), "reverse_proxy frontend:80") {
		t.Fatal("direct rendering did not normalize target")
	}
	path := filepath.Join(t.TempDir(), "config.yaml")
	cfg := Config{Path: path, Entries: []Entry{raw}}
	if err := cfg.Save(); err != nil {
		t.Fatal(err)
	}
	saved, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(saved), "frontend.test: frontend:80") {
		t.Fatalf("noncanonical saved target: %s", saved)
	}
	actualJSON, _, err := cfg.CaddyJSON()
	if err != nil {
		t.Fatal(err)
	}
	expectedJSON, _, err := (Config{Entries: []Entry{parsed}}).CaddyJSON()
	if err != nil {
		t.Fatal(err)
	}
	if string(actualJSON) != string(expectedJSON) || !strings.Contains(string(actualJSON), `"dial":"frontend:80"`) {
		t.Fatalf("unexpected adapted JSON: %s", actualJSON)
	}
	if cfg.Entries[0] != raw {
		t.Fatalf("config entry mutated: %#v", cfg.Entries[0])
	}
	cfg.Entries = []Entry{{Alias: "bad", Upstream: "8000"}}
	if err := cfg.Save(); err == nil {
		t.Fatal("numeric upstream accepted")
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(saved) {
		t.Fatal("invalid entry overwrote config")
	}
	if _, _, err := cfg.CaddyJSON(); err == nil {
		t.Fatal("numeric upstream accepted by CaddyJSON")
	}
}
