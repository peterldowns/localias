package server

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/hashicorp/mdns"

	_ "github.com/peterldowns/localias/pkg/caddymodules" // necessary caddy configuration
	"github.com/peterldowns/localias/pkg/config"
)

// Start and Stop are not safe to use in parallel environments
// but it's OK because that's not needed. They both mutate
// the `instance` package global.

var instance *Server //nolint:gochecknoglobals

func Start(cfg *config.Config) error {
	if instance != nil {
		if err := instance.Stop(); err != nil {
			return err
		}
	}
	instance = &Server{config: cfg}
	return instance.Start()
}

// Stop will stop the caddy server (if it is running).
func Stop() error {
	if instance == nil {
		return nil
	}
	if err := instance.Stop(); err != nil {
		return err
	}
	instance = nil
	return nil
}

type Server struct {
	config    *config.Config
	mdnserver *mdns.Server
}

func (s *Server) Start() error {
	// Start (or restart) the global Caddy service and load the current
	// configuration.
	cfgJSON, _, err := s.config.CaddyJSON()
	if err != nil {
		return err
	}
	if err := caddy.Load(cfgJSON, false); err != nil {
		return err
	}
	s.mdnserver, err = newMDNSServer(s.config.Entries)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop() error {
	if err := caddy.Stop(); err != nil {
		return err
	}
	if s.mdnserver != nil {
		if err := s.mdnserver.Shutdown(); err != nil {
			return err
		}
		s.mdnserver = nil
	}
	return nil
}
