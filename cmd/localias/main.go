package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"

	"github.com/peterldowns/localias/cmd/localias/root"
	"github.com/peterldowns/localias/cmd/localias/shared"
)

func main() {
	defer func() {
		switch t := recover().(type) {
		case error:
			onError(fmt.Errorf("panic: %w", t))
		case string:
			onError(fmt.Errorf("panic: %s", t))
		default:
			if t != nil {
				onError(fmt.Errorf("panic: %+v", t))
			}
		}
	}()
	if err := root.Command.Execute(); err != nil {
		onError(err)
	}
}

func onError(err error) {
	err = shared.ConvertErr(err)
	msg := fmt.Sprintf("error: %s", err)
	fmt.Fprintln(os.Stderr, color.New(color.FgRed, color.Italic).Sprintf(msg))
	os.Exit(1)
}
