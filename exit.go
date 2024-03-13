package cleo

import (
	"fmt"

	"github.com/markbates/plugins"
	"github.com/markbates/plugins/plugcmd"
)

type Exiter interface {
	// plugins.Stdioer
	plugins.Plugin
	Exit(code int) error
}

// Exit will print the usage information
// for the command to the given Stderr.
// If the cmd.ExitFn is set, it will be
// called with the given exit code.
// Otherwise, os.Exit will be called.
func Exit(cmd plugins.Stdioer, code int, err error) {
	if err == nil {
		return
	}

	if code == -1 && err == nil {
		plugcmd.Print(cmd.Stderr(), cmd)
		return
	}

	var p plugins.Plugin = cmd

	plugcmd.Print(cmd.Stderr(), p)
	fmt.Fprintf(cmd.Stderr(), "\nError: %s\n", err)

	if ex, ok := cmd.(Exiter); ok {
		ex.Exit(code)
	}

}
