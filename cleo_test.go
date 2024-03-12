package cleo

import (
	"context"
	"fmt"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/markbates/iox"
	"github.com/markbates/plugins"
	"github.com/markbates/plugins/plugcmd"
)

type stringPlug string

func (s stringPlug) PluginName() string {
	return string(s)
}

func newEcho(t testing.TB, name string) *echoPlug {
	t.Helper()
	e := &echoPlug{
		Name: name,
		IO:   iox.Discard(),
		FS:   fstest.MapFS{},
		Plugins: func() plugins.Plugins {
			return nil
		},
		// Desc: fmt.Sprintf("echo %s", name),

	}

	return e
}

var _ Exiter = &exiterPlug{}

type exiterPlug struct {
	Code int
}

func (e *exiterPlug) Exit(code int) error {
	if e == nil {
		return fmt.Errorf("nil exiter")
	}

	if code == -1 {
		return fmt.Errorf("exit code is -1")
	}
	e.Code = code
	return nil
}

func (e *exiterPlug) PluginName() string {
	return fmt.Sprintf("%T", e)
}

var _ plugcmd.Commander = &echoPlug{}
var _ plugcmd.Describer = &echoPlug{}
var _ plugcmd.SubCommander = &echoPlug{}
var _ plugins.FSSetable = &echoPlug{}
var _ plugins.Feeder = &echoPlug{}
var _ plugins.IOSetable = &echoPlug{}
var _ plugins.Needer = &echoPlug{}

type echoPlug struct {
	iox.IO

	Name    string
	FS      fs.FS
	Subs    []plugcmd.Commander
	Plugins plugins.FeederFn
}

func (e *echoPlug) WithPlugins(fn plugins.FeederFn) error {
	if e == nil {
		return fmt.Errorf("echoPlug is nil")
	}

	if fn == nil {
		return fmt.Errorf("fn is nil")
	}

	e.Plugins = fn

	return nil
}

func (e *echoPlug) PluginFeeder() plugins.FeederFn {
	fn := func() plugins.Plugins {
		return nil
	}

	if e == nil || e.Plugins == nil {
		return fn
	}

	return e.Plugins
}

func (e *echoPlug) Description() string {
	return fmt.Sprintf("echo %s", e.Name)
}

func (e *echoPlug) SubCommands() []plugcmd.Commander {
	return e.Subs
}

func (e *echoPlug) PluginName() string {
	return fmt.Sprintf("echo/%s", e.Name)
}

func (cmd *echoPlug) Main(ctx context.Context, pwd string, args []string) error {
	if cmd == nil {
		return fmt.Errorf("echoPlug is nil")
	}

	fmt.Fprint(cmd.IO.Stdout(), args)
	fmt.Fprint(cmd.IO.Stderr(), args)
	return nil
}

func (cmd *echoPlug) SetStdio(oi plugins.IO) error {
	if cmd == nil {
		return fmt.Errorf("echoPlug is nil")
	}

	cmd.IO = oi
	return nil
}

func (cmd *echoPlug) SetFileSystem(fs fs.FS) error {
	if cmd == nil {
		return fmt.Errorf("echoPlug is nil")
	}

	cmd.FS = fs
	return nil
}

func (cmd *echoPlug) FileSystem() (fs.FS, error) {
	if cmd == nil {
		return nil, fmt.Errorf("echoPlug is nil")
	}

	return cmd.FS, nil
}
