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

type Cmd struct {
	iox.IO // IO to be used by the command
	fs.FS  // FS to be used by the command
	sync.RWMutex

	Aliases  []string             // Aliases for the command
	Commands map[string]Commander // Sub commands for the command
	Feeder   plugins.Feeder       // Plugins for the command
	Name     string               // Name of the command

	Desc string // Description of the command

	ExitFn func(int) // ExitFn is used by the Exit method. Default: os.Exit
}

func (cmd *Cmd) Exit(code int) {
	if cmd == nil {
		return
	}

	if cmd.ExitFn == nil {
		return
	}

	cmd.ExitFn(code)
}

func (cmd *Cmd) Description() string {
	if cmd == nil {
		return ""
	}

	return cmd.Desc
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

	plugs := cmd.Feeder()

	res := make(plugins.Plugins, 0, len(plugs))
	for _, p := range plugs {
		if p == cmd {
			continue
		}
		res = append(res, p)
	}

	return res
}

// ScopedPlugins returns the plugins scoped to the command.
// If the plugins include the current command, it will be removed
// from the returned list.
func (cmd *Cmd) ScopedPlugins() plugins.Plugins {
	if cmd == nil {
		return nil
	}

	plugs := cmd.Plugins()

	return plugs
}

// SubCommands returns the sub-commands for the command.
func (cmd *Cmd) SubCommands() []Commander {
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

// Main is the main entry point for the command.
// NEEDS TO BE IMPLEMENTED
func (cmd *Cmd) Main(ctx context.Context, pwd string, args []string) error {
	return fmt.Errorf("not implemented")
}
