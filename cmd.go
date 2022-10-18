package cleo

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"sync"

	"github.com/markbates/iox"
	"github.com/markbates/plugins"
	"github.com/markbates/plugins/plugcmd"
)

type Commander = plugcmd.Commander

var _ FSSetable = &Cmd{}
var _ IOSetable = &Cmd{}
var _ iox.IOable = &Cmd{}
var _ plugcmd.Aliaser = &Cmd{}
var _ plugcmd.SubCommander = &Cmd{}
var _ plugins.FSable = &Cmd{}

type Cmd struct {
	iox.IO // IO to be used by the command
	fs.FS  // FS to be used by the command
	sync.RWMutex

	Name string // Name of the command

	Aliases []string       // Aliases for the command
	Feeder  plugins.Feeder // Plugins for the command

	ExitFn func(int) // ExitFn is used by the Exit method. Default: os.Exit
}

// Plugins will safely call the Feeder function
// if provided.
func (cmd *Cmd) Plugins() plugins.Plugins {
	if cmd == nil {
		return nil
	}

	cmd.RLock()
	defer cmd.RUnlock()

	if cmd.Feeder == nil {
		return nil
	}

	return cmd.Feeder()
}

// ScopedPlugins returns the plugins scoped to the command.
// If the plugins include the current command, it will be removed
// from the returned list.
func (cmd *Cmd) ScopedPlugins() plugins.Plugins {
	if cmd == nil {
		return nil
	}

	plugs := cmd.Plugins()

	res := make(plugins.Plugins, 0, len(plugs))
	for _, p := range plugs {
		if p == cmd {
			continue
		}
		res = append(res, p)
	}

	return res
}

// SubCommands returns the sub-commands for the command.
func (cmd *Cmd) SubCommands() plugins.Plugins {
	plugs := cmd.ScopedPlugins()
	if len(plugs) == 0 {
		return plugs
	}

	res := make(plugins.Plugins, 0, len(plugs))

	for _, p := range plugs {
		if _, ok := p.(Commander); ok {
			res = append(res, p)
		}
	}

	return res
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

	cmd.RLock()
	defer cmd.RUnlock()
	return cmd.Name
}

// CmdAliases returns the aliases for the command.
func (cmd *Cmd) CmdAliases() []string {
	if cmd == nil {
		return nil
	}

	cmd.RLock()
	defer cmd.RUnlock()
	return cmd.Aliases
}

// String returns a string representation of the command.
func (cmd *Cmd) String() string {
	if cmd == nil {
		return ""
	}

	b, _ := json.MarshalIndent(cmd, "", "  ")
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

	return json.MarshalIndent(m, "", "  ")
}
