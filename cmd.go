package cleo

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"sort"
	"sync"

	"github.com/markbates/iox"
	"github.com/markbates/plugins"
	"github.com/markbates/plugins/plugcmd"
)

type Commander = plugcmd.Commander

var _ Commander = &Cmd{}
var _ Exiter = &Cmd{}
var _ plugcmd.Describer = &Cmd{}
var _ plugins.FSSetable = &Cmd{}
var _ plugins.FSable = &Cmd{}
var _ plugins.Feeder = &Cmd{}
var _ plugins.IOSetable = &Cmd{}
var _ plugins.IOable = &Cmd{}
var _ plugins.Scoper = &Cmd{}

type Cmd struct {
	iox.IO // IO to be used by the command
	fs.FS  // FS to be used by the command

	Aliases  []string             // Aliases for the command
	Commands map[string]Commander // Sub commands for the command
	Feeder   plugins.FeederFn     // Plugins for the command
	Name     string               // Name of the command

	Desc string // Description of the command

	ExitFn func(int) error // ExitFn is used by the Exit method.

	mu sync.RWMutex
}

func (cmd *Cmd) Exit(code int) error {
	if cmd == nil {
		return fmt.Errorf("nil command")
	}

	plugs := cmd.ScopedPlugins()

	exiters := plugins.ByType[Exiter](plugs)
	for _, ex := range exiters {
		if err := ex.Exit(code); err != nil {
			return err
		}
	}

	cmd.mu.RLock()
	fn := cmd.ExitFn
	cmd.mu.RUnlock()

	if fn == nil {
		return nil
	}

	return fn(code)
}

func (cmd *Cmd) Description() string {
	if cmd == nil {
		return ""
	}

	cmd.mu.RLock()
	defer cmd.mu.RUnlock()

	return cmd.Desc
}

// Plugins will provider a single FeederFn
// that will return all of the plugins that
// are available to the command.
func (cmd *Cmd) PluginFeeder() plugins.FeederFn {
	fn := func() plugins.Plugins {
		return nil
	}

	if cmd == nil {
		return fn
	}

	cmd.mu.RLock()
	defer cmd.mu.RUnlock()

	if cmd.Feeder == nil {
		return fn
	}

	return func() plugins.Plugins {
		var plugs plugins.Plugins
		if cmd.Feeder != nil {
			plugs = cmd.Feeder()
		}

		for _, p := range plugs {
			if pf, ok := p.(plugins.Feeder); ok {
				fn := pf.PluginFeeder()
				plugs = append(plugs, fn()...)
			}
		}

		return plugs
	}
}

// ScopedPlugins returns the plugins scoped to the command.
// If the plugins include the current command, it will be removed
// from the returned list.
func (cmd *Cmd) ScopedPlugins() plugins.Plugins {
	return cmd.PluginFeeder()()
}

// SubCommands returns the sub-commands for the command.
func (cmd *Cmd) SubCommands() []Commander {
	if cmd == nil {
		return nil
	}

	cmd.mu.RLock()
	defer cmd.mu.RUnlock()

	cmds := make([]Commander, 0, len(cmd.Commands))

	keys := make([]string, 0, len(cmd.Commands))

	for k := range cmd.Commands {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		cmds = append(cmds, cmd.Commands[k])
	}

	return cmds
}

// PluginName returns name of the plugin.
func (cmd *Cmd) PluginName() string {
	name := "?"
	if cmd != nil {
		name = cmd.CmdName()
	}

	return fmt.Sprintf("%T (%s)", cmd, name)

}

// CmdName returns the name of the command.
func (cmd *Cmd) CmdName() string {
	if cmd == nil {
		return ""
	}

	cmd.mu.RLock()
	defer cmd.mu.RUnlock()
	return cmd.Name
}

// CmdAliases returns the aliases for the command.
func (cmd *Cmd) CmdAliases() []string {
	if cmd == nil {
		return nil
	}

	cmd.mu.RLock()
	defer cmd.mu.RUnlock()
	return cmd.Aliases
}

// String returns a string representation of the command.
func (cmd *Cmd) String() string {
	if cmd == nil {
		return ""
	}

	b, _ := json.Marshal(cmd)
	return string(b)
}

// MarshalJSON returns a JSON representation of the command.
func (cmd *Cmd) MarshalJSON() ([]byte, error) {
	if cmd == nil {
		return nil, fmt.Errorf("nil command")
	}

	plugs := cmd.ScopedPlugins()

	m := map[string]any{
		"aliases": cmd.Aliases,
		"name":    cmd.Name,
		"stdio":   cmd.Stdio(),
		"plugins": plugs,
	}

	return json.Marshal(m)
}

// Main is the main entry point for the command.
// NEEDS TO BE IMPLEMENTED
func (cmd *Cmd) Main(ctx context.Context, pwd string, args []string) error {
	return fmt.Errorf("not implemented")
}
