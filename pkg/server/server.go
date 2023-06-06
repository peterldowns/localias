package server

import (
	"fmt"
	"net"
	"os"

	"github.com/caddyserver/caddy/v2"
	"github.com/hashicorp/mdns"

	_ "github.com/peterldowns/localias/pkg/caddymodules" // necessary caddy configuration
	"github.com/peterldowns/localias/pkg/config"
)

// Start will start the caddy server (if it hasn't been started already).
func Start(cfg *config.Config) error {
	cfgJSON, _, err := cfg.CaddyJSON()
	if err != nil {
		return err
	}
	if err := caddy.Load(cfgJSON, false); err != nil {
		return err
	}

	// Setup our service export
	host, err := os.Hostname()
	if err != nil {
		return err
	}
	info := []string{"frontend.local/"}
	ips, err := net.LookupIP(host + ".local")
	if err != nil {
		return fmt.Errorf("failed to find local IP: %w", err)
	}
	service, err := mdns.NewMDNSService(host, "localias._tcp", "", "frontend.local.", 3000, ips, info)
	if err != nil {
		return err
	}
	// Create the mDNS server, defer shutdown
	_, err = mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		return err
	}
	service, err = mdns.NewMDNSService(host, "localias._tcp", "", "backend.local.", 3000, ips, info)
	if err != nil {
		return err
	}
	// Create the mDNS server, defer shutdown
	_, err = mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		return err
	}
	select {}
	// defer server.Shutdown()
	return nil
}

// Stop will stop the caddy server (if it is running).
func Stop() error {
	return caddy.Stop()
}
