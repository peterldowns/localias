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
	// var msg string
	// var le shared.LocaliasError
	// if errors.As(err, &le) {
	// 	msg = fmt.Sprintf("error(%s): %s", le.Code(), le.Error())
	// } else {
	// 	msg = fmt.Sprintf("error: %s", err)
	// }
	msg := fmt.Sprintf("error: %s", err)
	fmt.Fprintln(os.Stderr, color.New(color.FgRed, color.Italic).Sprintf(msg))
	os.Exit(1)
}
