package cleo

import (
	"errors"
	"fmt"
	"os"

	"github.com/markbates/plugins"
	"github.com/markbates/plugins/plugcmd"
)

// Exit will print the usage information
// for the command to the given Stderr.
// If the cmd.ExitFn is set, it will be
// called with the given exit code.
// Otherwise, os.Exit will be called.
func (cmd *Cmd) Exit(i int, err error) {
	if err == nil {
		return
	}

	var p plugins.Plugin = cmd

	var e plugins.Error
	if errors.As(err, &e) {
		p = e.Plugin
	}

	plugcmd.Print(cmd.Stderr(), p)
	fmt.Fprintf(cmd.Stderr(), "\nError: %s\n", err)

	if cmd.ExitFn == nil {
		os.Exit(i)
	}

	cmd.ExitFn(i)
}
