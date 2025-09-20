package cleo

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"sort"
	"strings"
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

	logger *slog.Logger // Structured logger
	mu     sync.RWMutex
}

// NewCmd creates a new command with the given name and options.
func NewCmd(name string, opts ...Option) (*Cmd, error) {
	if name == "" {
		return nil, ErrEmptyCommandName
	}

	cmd := &Cmd{
		Name:     name,
		Commands: make(map[string]Commander),
		logger:   slog.Default(),
	}

	for _, opt := range opts {
		if err := opt(cmd); err != nil {
			return nil, fmt.Errorf("applying option: %w", err)
		}
	}

	return cmd, nil
}

// Logger returns the logger for the command.
func (cmd *Cmd) Logger() *slog.Logger {
	if cmd == nil {
		return slog.Default()
	}
	cmd.mu.RLock()
	defer cmd.mu.RUnlock()
	if cmd.logger == nil {
		return slog.Default()
	}
	return cmd.logger
}

func (cmd *Cmd) Exit(code int) error {
	return cmd.ExitWithContext(context.Background(), code)
}

// ExitWithContext exits the command with the given code, using the provided context.
func (cmd *Cmd) ExitWithContext(ctx context.Context, code int) error {
	if cmd == nil {
		return ErrNilCommand
	}

	logger := cmd.Logger()
	logger.InfoContext(ctx, "Command exiting", "code", code, "name", cmd.CmdName())

	plugs := cmd.ScopedPlugins()

	exiters := plugins.ByType[Exiter](plugs)
	for _, ex := range exiters {
		if err := ex.Exit(code); err != nil {
			logger.ErrorContext(ctx, "Plugin exit failed", "error", err, "plugin", ex.PluginName())
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
	feeder := cmd.PluginFeeder()
	if feeder == nil {
		return nil
	}

	// Use pool for temporary plugin collection
	tempPlugs := getPluginSlice()
	defer putPluginSlice(tempPlugs)

	plugs := feeder()
	if len(plugs) == 0 {
		return nil
	}

	// Copy to result slice
	result := make(plugins.Plugins, len(plugs))
	copy(result, plugs)

	return result
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
		return nil, ErrNilCommand
	}

	var buf strings.Builder
	buf.WriteString(`{`)

	// Write aliases
	buf.WriteString(`"aliases":`)
	aliasesJSON, err := json.Marshal(cmd.Aliases)
	if err != nil {
		return nil, fmt.Errorf("marshaling aliases: %w", err)
	}
	buf.Write(aliasesJSON)

	// Write name
	buf.WriteString(`,"name":`)
	nameJSON, err := json.Marshal(cmd.Name)
	if err != nil {
		return nil, fmt.Errorf("marshaling name: %w", err)
	}
	buf.Write(nameJSON)

	// Write stdio
	buf.WriteString(`,"stdio":`)
	stdioJSON, err := json.Marshal(cmd.Stdio())
	if err != nil {
		return nil, fmt.Errorf("marshaling stdio: %w", err)
	}
	buf.Write(stdioJSON)

	// Write plugins
	buf.WriteString(`,"plugins":`)
	plugs := cmd.ScopedPlugins()
	plugsJSON, err := json.Marshal(plugs)
	if err != nil {
		return nil, fmt.Errorf("marshaling plugins: %w", err)
	}
	buf.Write(plugsJSON)

	buf.WriteString(`}`)
	return []byte(buf.String()), nil
}

// Main is the main entry point for the command.
// NEEDS TO BE IMPLEMENTED
func (cmd *Cmd) Main(ctx context.Context, pwd string, args []string) error {
	return cmd.MainWithContext(ctx, pwd, args)
}

// MainWithContext is the main entry point for the command with context support.
// NEEDS TO BE IMPLEMENTED
func (cmd *Cmd) MainWithContext(ctx context.Context, pwd string, args []string) error {
	if cmd == nil {
		return ErrNilCommand
	}

	logger := cmd.Logger()
	logger.InfoContext(ctx, "Command main called",
		"name", cmd.CmdName(),
		"pwd", pwd,
		"args", args)

	return fmt.Errorf("not implemented")
}
