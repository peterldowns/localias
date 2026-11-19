package config

import (
	"fmt"
	"net"
	"net/netip"
	"strconv"
	"strings"

	"github.com/miekg/dns"
)

const (
	// TCP ports are unsigned 16-bit values; zero is not a dialable service port.
	minUpstreamPort = 1
	maxUpstreamPort = 65535
	portBitSize     = 16
	decimalRadix    = 10
	// Maximum textual DNS name length, excluding the optional root dot.
	maxDNSNameLen = 253
)

// ParseTarget parses a local port or an explicit HTTP upstream address.
func ParseTarget(value any) (Entry, error) {
	switch v := value.(type) {
	case int:
		if v < minUpstreamPort || v > maxUpstreamPort {
			return Entry{}, fmt.Errorf("port must be between %d and %d", minUpstreamPort, maxUpstreamPort)
		}
		return Entry{Port: v}, nil
	case string:
		if v == "" {
			return Entry{}, fmt.Errorf("target must not be empty")
		}
		// Numeric strings represent local ports and normalize to the integer form.
		if digits(v) {
			p, err := parsePort(v)
			if err != nil {
				return Entry{}, err
			}
			return Entry{Port: p}, nil
		}
		// Validate syntax without resolving DNS: services may not exist until runtime.
		host, port, err := net.SplitHostPort(v)
		if err != nil || host == "" {
			return Entry{}, fmt.Errorf("target must be a port or host:port")
		}
		p, err := parsePort(port)
		if err != nil {
			return Entry{}, err
		}
		// Brackets are reserved for IPv6 literals; zones are interface-specific.
		if strings.HasPrefix(v, "[") {
			ip, err := netip.ParseAddr(host)
			if err != nil || !ip.Is6() || ip.Zone() != "" {
				return Entry{}, fmt.Errorf("invalid IPv6 host %q", host)
			}
			host = ip.String()
		} else if !validHost(host) {
			return Entry{}, fmt.Errorf("invalid upstream host %q", host)
		}
		return Entry{Upstream: net.JoinHostPort(host, strconv.Itoa(p))}, nil
	default:
		return Entry{}, fmt.Errorf("target must be an integer or string, got %T", value)
	}
}

func digits(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func parsePort(s string) (int, error) {
	p, err := strconv.ParseUint(s, decimalRadix, portBitSize)
	if err != nil || p < minUpstreamPort {
		return 0, fmt.Errorf("port must be between %d and %d", minUpstreamPort, maxUpstreamPort)
	}
	return int(p), nil
}

func validHost(host string) bool {
	if ip, err := netip.ParseAddr(host); err == nil {
		return ip.Is4()
	}
	// The DNS helper validates labels without resolving the name.
	if _, ok := dns.IsDomainName(host); !ok {
		return false
	}
	// An absolute DNS name may have a trailing dot, which is not a label.
	host = strings.TrimSuffix(host, ".")
	// Valid IPv4 literals return early above; reject any non-IPv4 host whose
	// rightmost label is all digits (RFC 1123 / RFC 3696), including malformed
	// or trailing-dot IPv4 addresses.
	parts := strings.Split(host, ".")
	if digits(parts[len(parts)-1]) {
		return false
	}
	// IsDomainName accepts names beyond the DNS wire limit.
	if len(host) > maxDNSNameLen {
		return false
	}
	for _, label := range parts {
		if len(label) == 0 || label[0] == '-' || label[len(label)-1] == '-' {
			return false
		}
		// Hosts are interpolated into a Caddyfile, so exclude whitespace and
		// configuration delimiters. Allow underscores used by container services.
		for _, c := range label {
			if (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') && (c < '0' || c > '9') && c != '-' && c != '_' {
				return false
			}
		}
	}
	return true
}

// Target returns the target in its user-facing form.
func (entry Entry) Target() string {
	if normalized, err := entry.NormalizeTarget(); err == nil {
		entry = normalized
	}
	if entry.Upstream != "" {
		return entry.Upstream
	}
	return strconv.Itoa(entry.Port)
}

// NormalizeTarget validates and canonicalizes the target while preserving the alias.
func (entry Entry) NormalizeTarget() (Entry, error) {
	// Equal presence flags mean both targets are set or neither is set.
	if (entry.Port != 0) == (entry.Upstream != "") {
		return Entry{}, fmt.Errorf("exactly one port or upstream must be set")
	}
	var parsed Entry
	var err error
	if entry.Upstream != "" {
		parsed, err = ParseTarget(entry.Upstream)
		if err != nil {
			return Entry{}, err
		}
		if parsed.Upstream == "" {
			return Entry{}, fmt.Errorf("upstream must be a host:port address")
		}
	} else {
		parsed, err = ParseTarget(entry.Port)
		if err != nil {
			return Entry{}, err
		}
	}
	parsed.Alias = entry.Alias
	return parsed, nil
}

// ValidateTarget ensures exactly one valid target is set.
func (entry Entry) ValidateTarget() error {
	_, err := entry.NormalizeTarget()
	return err
}

func (entry Entry) invalidTargetError(err error) error {
	switch {
	case entry.Port == 0 && entry.Upstream == "":
		return fmt.Errorf("alias %q: invalid target: %w", entry.Alias, err)
	case entry.Port != 0 && entry.Upstream != "":
		return fmt.Errorf("alias %q: invalid target (port=%d, upstream=%q): %w", entry.Alias, entry.Port, entry.Upstream, err)
	case entry.Upstream != "":
		return fmt.Errorf("alias %q: invalid target %q: %w", entry.Alias, entry.Upstream, err)
	default:
		return fmt.Errorf("alias %q: invalid target %d: %w", entry.Alias, entry.Port, err)
	}
}

func (entry Entry) rawUpstreamAddress() string {
	if entry.Upstream != "" {
		return entry.Upstream
	}
	return net.JoinHostPort("localhost", strconv.Itoa(entry.Port))
}

// UpstreamAddress returns the address used by the reverse proxy.
func (entry Entry) UpstreamAddress() string {
	if normalized, err := entry.NormalizeTarget(); err == nil {
		return normalized.rawUpstreamAddress()
	}
	return entry.rawUpstreamAddress()
}
