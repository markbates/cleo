package cleo

import (
	"context"

	"github.com/markbates/plugins"
)

// Init will initialize the command.
// It should be called before the command is used.
func (cmd *Cmd) Init() error {
	return cmd.InitWithContext(context.Background())
}

// InitWithContext will initialize the command with the given context.
// It should be called before the command is used.
func (cmd *Cmd) InitWithContext(ctx context.Context) error {
	if cmd == nil {
		return ErrNilCommand
	}

	logger := cmd.Logger()
	logger.InfoContext(ctx, "Initializing command", "name", cmd.CmdName())

	cmd.mu.Lock()
	cab := cmd.FS
	cmd.mu.Unlock()

	plugFn := cmd.PluginFeeder()
	plugs := plugFn()

	logger.InfoContext(ctx, "Found plugins", "count", len(plugs))

	// plugins.FSSetable
	fss := plugins.ByType[plugins.FSSetable](plugs)
	for _, fs := range fss {
		if err := fs.SetFileSystem(cab); err != nil {
			logger.ErrorContext(ctx, "Failed to set filesystem for plugin", "error", err)
			return err
		}
	}

	// plugins.Needer
	needs := plugins.ByType[plugins.Needer](plugs)
	for _, n := range needs {
		if err := n.WithPlugins(plugFn); err != nil {
			logger.ErrorContext(ctx, "Failed to set plugins for needer", "error", err)
			return err
		}
	}

	// plugins.IOSetable
	ios := plugins.ByType[plugins.IOSetable](plugs)
	for _, io := range ios {
		if err := io.SetStdio(cmd.IO); err != nil {
			logger.ErrorContext(ctx, "Failed to set stdio for plugin", "error", err)
			return err
		}
	}

	logger.InfoContext(ctx, "Command initialization completed", "name", cmd.CmdName())
	return nil
}
