package server

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/hashicorp/mdns"
	"github.com/miekg/dns"

	"github.com/peterldowns/localias/pkg/config"
)

// multiservice implements the mdns.Zone interface, and will respond to dns
// questions by fanning them out to multiple services.
type multiservice []*mdns.MDNSService

func (m multiservice) Records(q dns.Question) []dns.RR {
	var records []dns.RR
	for _, s := range m {
		records = append(records, s.Records(q)...)
	}
	// (This is a good point to add logging/tracing in order to debug mDNS)
	return records
}

// newMDNSServer creates and starts a mDNS server if there are any aliases
// ending in ".local". While this is running, other devices on the same wifi
// network will be able to visit these aliases.
func newMDNSServer(entries []config.Entry) (*mdns.Server, error) {
	var localEntries []config.Entry
	for _, entry := range entries {
		if isLocal(entry) {
			localEntries = append(localEntries, entry)
		}
	}
	if localEntries == nil {
		return nil, nil
	}
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	// To start the mDNS server, have to look up the IP of the host machine that
	// it will run on. If the hostname doesn't have a `.local` suffix, we need
	// to add one for some reason based on the example mDNS code I've found.
	// If the hostname already has a `.local` suffix, we should keep it.
	baseIPs, err := getIPAddressesForMDNS()
	if err != nil {
		return nil, err
	}
	if baseIPs == nil {
		return nil, fmt.Errorf("unable to serve mDNS: could not determine a local IP for host %s", hostname)
	}
	var ms multiservice
	for _, entry := range localEntries {
		ehost := entry.Host()
		service, err := mdns.NewMDNSService(
			// Necessary to escape the periods in an instance name for some reason,
			// verified by testing with Discovery.app and with
			// `dns-sd -B _http._tcp local`
			//
			// 		foo\.local
			//
			strings.ReplaceAll(ehost, ".", "\\."),
			// Use _http for both _https and _http services, since the _https
			// services will have a _http redirect anyway.
			//
			// http://www.dns-sd.org/ServiceTypes.html
			"_http._tcp",
			// The default value is "local." and seems like it shouldn't ever be
			// anything else.
			"local.",
			// The hostname, including the TLD ("local") and a trailing ".", is
			// what is used to actually answer mDNS queries.
			//
			// 		foo.local.
			//
			ehost+".",
			// Instead of using the service port directly we proxy through
			// Caddy, so we use either port 443 (for secure aliases) or 80 (for
			// insecure aliases).
			//
			// 		443
			//
			caddyPort(entry),
			// Use the IP addresses we looked up earlier for the host machine as
			// the answer to "which IPs can be used to access this alias/host",
			// since this machine is where Caddy and where the service is
			// actually running.
			baseIPs,
			// Just for fun, include a TXT record giving Localias credit.
			[]string{ehost + " @ " + hostname + " via localias"},
			// nil,
		)
		if err != nil {
			return nil, err
		}
		fmt.Printf("mDNS: serving %s (%s)\n", entry.Host(), ipStrings(baseIPs))
		ms = append(ms, service)
	}
	return mdns.NewServer(&mdns.Config{Zone: ms})
}

func isLocal(entry config.Entry) bool {
	return strings.HasSuffix(entry.Host(), ".local")
}

func caddyPort(entry config.Entry) int {
	a, _ := httpcaddyfile.ParseAddress(entry.Alias)
	if a.Scheme == "" {
		a.Scheme = "https"
	}
	if a.Scheme == "https" {
		return 443
	}
	return 80
}

func ensureSuffix(s, suffix string) string {
	return strings.TrimSuffix(s, suffix) + suffix
}

func getIPAddressesForMDNS() ([]net.IP, error) {
	// Get a list of all network interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var ipAddresses []net.IP
	for _, iface := range interfaces {
		// Get a list of all unicast IP addresses associated with the interface
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			// Exclude local loopback addresses and V6 addresses,
			if ip == nil || ip.IsLoopback() || ip.To4() == nil {
				continue
			}
			ipAddresses = append(ipAddresses, ip)
		}
	}
	// On MacOS Sonoma 14.4.1, if you don't include the IPV6 link-local IP in
	// the broadcast, DNS lookups will hang for like 5-6 seconds each time, even
	// with the IPV4 address in /etc/hosts. I can't explain it but including
	// fe80::1 makes the lookups instant.
	if len(ipAddresses) != 0 {
		ipAddresses = append(ipAddresses, net.ParseIP("::1"), net.ParseIP("127.0.0.1"))
	}
	return ipAddresses, nil
}

func ipStrings(ips []net.IP) string {
	var ipStrings []string
	for _, ip := range ips {
		ipStrings = append(ipStrings, ip.String())
	}
	return strings.Join(ipStrings, ",")
}
