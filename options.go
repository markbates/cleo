package cleo

import (
	"io/fs"
	"log/slog"

	"github.com/markbates/iox"
	"github.com/markbates/plugins"
)

// Option represents a functional option for configuring a Cmd.
type Option func(*Cmd) error

// WithStdio sets the I/O for the command.
func WithStdio(io iox.IO) Option {
	return func(cmd *Cmd) error {
		if cmd == nil {
			return ErrNilCommand
		}
		cmd.mu.Lock()
		cmd.IO = io
		cmd.mu.Unlock()
		return nil
	}
}

// WithFS sets the filesystem for the command.
func WithFS(filesystem fs.FS) Option {
	return func(cmd *Cmd) error {
		if cmd == nil {
			return ErrNilCommand
		}
		if filesystem == nil {
			return ErrNilFS
		}
		cmd.mu.Lock()
		cmd.FS = filesystem
		cmd.mu.Unlock()
		return nil
	}
}

// WithAliases sets the aliases for the command.
func WithAliases(aliases ...string) Option {
	return func(cmd *Cmd) error {
		if cmd == nil {
			return ErrNilCommand
		}
		cmd.mu.Lock()
		cmd.Aliases = make([]string, len(aliases))
		copy(cmd.Aliases, aliases)
		cmd.mu.Unlock()
		return nil
	}
}

// WithDescription sets the description for the command.
func WithDescription(desc string) Option {
	return func(cmd *Cmd) error {
		if cmd == nil {
			return ErrNilCommand
		}
		cmd.mu.Lock()
		cmd.Desc = desc
		cmd.mu.Unlock()
		return nil
	}
}

// WithPluginFeeder sets the plugin feeder for the command.
func WithPluginFeeder(feeder plugins.FeederFn) Option {
	return func(cmd *Cmd) error {
		if cmd == nil {
			return ErrNilCommand
		}
		cmd.mu.Lock()
		cmd.Feeder = feeder
		cmd.mu.Unlock()
		return nil
	}
}

// WithExitFunction sets the exit function for the command.
func WithExitFunction(exitFn func(int) error) Option {
	return func(cmd *Cmd) error {
		if cmd == nil {
			return ErrNilCommand
		}
		cmd.mu.Lock()
		cmd.ExitFn = exitFn
		cmd.mu.Unlock()
		return nil
	}
}

// WithLogger sets the logger for the command.
func WithLogger(logger *slog.Logger) Option {
	return func(cmd *Cmd) error {
		if cmd == nil {
			return ErrNilCommand
		}
		cmd.mu.Lock()
		cmd.logger = logger
		cmd.mu.Unlock()
		return nil
	}
}

// WithSubCommand adds a sub-command to the command.
func WithSubCommand(name string, subcmd Commander) Option {
	return func(cmd *Cmd) error {
		if cmd == nil {
			return ErrNilCommand
		}
		if name == "" {
			return ErrEmptyCommandName
		}
		if subcmd == nil {
			return ErrNilSubCommand
		}
		cmd.mu.Lock()
		if cmd.Commands == nil {
			cmd.Commands = make(map[string]Commander)
		}
		cmd.Commands[name] = subcmd
		cmd.mu.Unlock()
		return nil
	}
}