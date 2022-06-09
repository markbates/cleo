package cleo

import (
	"context"
	"fmt"
	"testing"
)

func newEcho(t testing.TB) *echo {
	t.Helper()
	return &echo{
		Cmd: &Cmd{},
	}
}

type echo struct {
	*Cmd
}

func (cmd *echo) Main(ctx context.Context, pwd string, args []string) error {
	if cmd.Cmd == nil {
		cmd.Cmd = &Cmd{}
	}

	fmt.Fprint(cmd.Stdout(), args)
	fmt.Fprint(cmd.Stderr(), args)
	return nil
}
