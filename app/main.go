package main

import (
	"C"
	"encoding/json"

	"github.com/peterldowns/localias/pkg/config"
)

import (
	"github.com/peterldowns/localias/pkg/daemon"
	"github.com/peterldowns/localias/pkg/hostctl"
)

//export config_open
func config_open() *C.char {
	cfg, err := config.Load(nil)
	if err != nil {
		return C.CString(err.Error())
	}
	bytes, err := json.Marshal(cfg)
	if err != nil {
		return C.CString(err.Error())
	}
	return C.CString(string(bytes))
}

//export config_save
func config_save(cfgjson *C.char) *C.char {
	var cfg config.Config
	bytes := []byte(C.GoString(cfgjson))
	if err := json.Unmarshal(bytes, &cfg); err != nil {
		return C.CString(err.Error())
	}
	if err := cfg.Save(); err != nil {
		return C.CString(err.Error())
	}
	return nil
}

//export daemon_start
func daemon_start() *C.char {
	cfg, _ := config.Load(nil)
	hctl := hostctl.DefaultController()
	if err := daemon.Run(hctl, cfg); err != nil {
		return C.CString(err.Error())
	}
	return nil
}

//export daemon_stop
func daemon_stop() *C.char {
	if err := daemon.Stop(); err != nil {
		return C.CString(err.Error())
	}
	return nil
}

// This entry point is somehow necessary for CGo.
func main() {}
