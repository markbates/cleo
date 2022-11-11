package cleo

import (
	"errors"
	"fmt"

	"github.com/markbates/plugins"
	"github.com/markbates/plugins/plugcmd"
)

type Exiter interface {
	plugins.Stdioer
	Exit(int)
}

// Exit will print the usage information
// for the command to the given Stderr.
// If the cmd.ExitFn is set, it will be
// called with the given exit code.
// Otherwise, os.Exit will be called.
func Exit(cmd plugins.Stdioer, i int, err error) {
	if err == nil {
		return
	}

	if i == -1 && err == nil {
		plugcmd.Print(cmd.Stderr(), cmd)
		return
	}

	var p plugins.Plugin = cmd

	var e plugins.Error
	if errors.As(err, &e) {
		p = e.Plugin
	}

	plugcmd.Print(cmd.Stderr(), p)
	fmt.Fprintf(cmd.Stderr(), "\nError: %s\n", err)

	if ex, ok := cmd.(Exiter); ok {
		ex.Exit(i)
	}

}
