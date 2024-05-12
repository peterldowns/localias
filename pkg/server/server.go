package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/caddyserver/caddy/v2"
	_ "github.com/caddyserver/caddy/v2/modules/standard"
	"github.com/fatih/color"
	"github.com/hashicorp/mdns"

	"github.com/peterldowns/localias/pkg/config"
)

func WaitForExitSignal() {
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
	<-quitChannel
}

type Server struct {
	Config     *config.Config
	MDNSServer *mdns.Server
}

func (s *Server) Start() error {
	if err := s.StartCaddy(); err != nil {
		return err
	}
	return s.StartMDNS()
}

func (s *Server) StartCaddy() error {
	// Start (or restart) the global Caddy service and load the current
	// configuration.
	cfgJSON, _, err := s.Config.CaddyJSON()
	if err != nil {
		return err
	}
	return caddy.Load(cfgJSON, false)
}

func (s *Server) StartMDNS() error {
	var err error
	s.MDNSServer, err = newMDNSServer(s.Config.Entries)
	if err != nil {
		warn := color.New(color.FgYellow, color.Italic)
		fmt.Fprintln(os.Stderr, warn.Sprintf("failed to start mDNS server:"))
		fmt.Fprintln(os.Stderr, warn.Sprintf(err.Error()))
	}
	return nil
}

func (s *Server) Stop() error {
	if err := caddy.Stop(); err != nil {
		return err
	}
	if s.MDNSServer != nil {
		if err := s.MDNSServer.Shutdown(); err != nil {
			return err
		}
		s.MDNSServer = nil
	}
	return nil
}
