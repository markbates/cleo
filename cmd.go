package cleo

import (
	"context"
	"io/fs"
	"sync"
)

type Commander interface {
	Main(ctx context.Context, pwd string, args []string) error
}

var _ Commander = &Cmd{}
var _ FSable = &Cmd{}
var _ IOSetable = &Cmd{}
var _ IOable = &Cmd{}
var _ SetFSable = &Cmd{}

type Cmd struct {
	IO
	fs.FS

	subs map[string]Commander
	sync.RWMutex
}

func (cmd *Cmd) Add(route string, c Commander) {
	if cmd == nil || c == nil {
		return
	}

	cmd.Lock()
	defer cmd.Unlock()

	if cmd.subs == nil {
		cmd.subs = map[string]Commander{}
	}

	cmd.subs[route] = c
}

func (cmd *Cmd) FileSystem() fs.FS {
	if cmd == nil {
		return nil
	}

	cmd.RLock()
	defer cmd.RUnlock()
	return cmd.FS
}

func (cmd *Cmd) SetFileSystem(cab fs.FS) {
	if cmd == nil {
		return
	}

	cmd.Lock()
	defer cmd.Unlock()

	cmd.FS = cab
}

func (cmd *Cmd) Stdio() IO {
	if cmd == nil {
		return IO{}
	}

	cmd.RLock()
	defer cmd.RUnlock()
	return cmd.IO
}

func (cmd *Cmd) SetStdio(oi IO) {
	if cmd == nil {
		return
	}

	cmd.Lock()
	defer cmd.Unlock()
	cmd.IO = oi
}

func (cmd *Cmd) Main(ctx context.Context, pwd string, args []string) error {
	if cmd == nil || cmd.subs == nil || len(cmd.subs) == 0 {
		return ErrNoCommands
	}

	if len(args) == 0 {
		return ErrNoCommand
	}

	cmd.RLock()

	c, ok := cmd.subs[args[0]]
	if !ok {
		cmd.RUnlock()
		return ErrUnknownCommand(args[0])
	}
	cmd.RLock()

	if ioc, ok := c.(IOSetable); ok {
		ioc.SetStdio(cmd.IO)
	}

	if cab, ok := c.(SetFSable); ok {
		cab.SetFileSystem(cmd.FS)
	}

	return c.Main(ctx, pwd, args[1:])
}
