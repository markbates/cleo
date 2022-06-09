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
	mu   sync.RWMutex
}

func (cmd *Cmd) Add(route string, c Commander) {
	if cmd == nil || c == nil {
		return
	}

	cmd.mu.Lock()
	defer cmd.mu.Unlock()

	if cmd.subs == nil {
		cmd.subs = map[string]Commander{}
	}

	cmd.subs[route] = c
}

func (cmd *Cmd) FileSystem() fs.FS {
	if cmd == nil {
		return nil
	}

	cmd.mu.RLock()
	defer cmd.mu.RUnlock()
	return cmd.FS
}

func (cmd *Cmd) SetFileSystem(cab fs.FS) {
	if cmd == nil {
		return
	}

	cmd.mu.Lock()
	defer cmd.mu.Unlock()

	cmd.FS = cab
}

func (cmd *Cmd) Stdio() IO {
	if cmd == nil {
		return IO{}
	}

	cmd.mu.RLock()
	defer cmd.mu.RUnlock()
	return cmd.IO
}

func (cmd *Cmd) SetStdio(oi IO) {
	if cmd == nil {
		return
	}

	cmd.mu.Lock()
	defer cmd.mu.Unlock()
	cmd.IO = oi
}

func (cmd *Cmd) Main(ctx context.Context, pwd string, args []string) error {
	if cmd == nil || cmd.subs == nil || len(cmd.subs) == 0 {
		return ErrNoCommands
	}

	if len(args) == 0 {
		return ErrNoCommand
	}

	cmd.mu.RLock()

	c, ok := cmd.subs[args[0]]
	if !ok {
		cmd.mu.RUnlock()
		return ErrUnknownCommand(args[0])
	}
	cmd.mu.RLock()

	if ioc, ok := c.(IOSetable); ok {
		ioc.SetStdio(cmd.IO)
	}

	if cab, ok := c.(SetFSable); ok {
		cab.SetFileSystem(cmd.FS)
	}

	return c.Main(ctx, pwd, args[1:])
}
