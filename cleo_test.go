package cleo

import (
	"context"
	"fmt"
	"testing"

	"github.com/markbates/iox"
	"github.com/markbates/plugins"
)

type String string

func (s String) PluginName() string {
	return string(s)
}

func newEcho(t testing.TB, name string) *echoPlug {
	t.Helper()
	e := &echoPlug{
		Cmd: &Cmd{
			Name: name,
			IO:   iox.Discard(),
		},
	}

	return e
}

type echoPlug struct {
	*Cmd
}

func (e *echoPlug) PluginName() string {
	return "echo"
}

func (cmd *echoPlug) Main(ctx context.Context, pwd string, args []string) error {
	if cmd.Cmd == nil {
		cmd.Cmd = &Cmd{}
	}

	fmt.Fprint(cmd.Stdout(), args)
	fmt.Fprint(cmd.Stderr(), args)
	return nil
}

func (cmd *echoPlug) SetStdio(oi plugins.IO) {
	cmd.IO = oi
}
