package cleo

import (
	"context"
	"log/slog"

	"github.com/markbates/plugins/plugcmd"
)

// Core interfaces for better separation of concerns

// Namer provides command naming functionality.
type Namer interface {
	CmdName() string
}

// Aliaser provides command aliasing functionality.
type Aliaser interface {
	CmdAliases() []string
}

// Describer provides command description functionality.
type Describer interface {
	Description() string
}

// Logger provides structured logging functionality.
type Logger interface {
	Logger() *slog.Logger
}

// ContextualCommander extends the base Commander with context support.
type ContextualCommander interface {
	plugcmd.Commander
	MainWithContext(ctx context.Context, pwd string, args []string) error
}

// ContextualExiter provides context-aware exit functionality.
type ContextualExiter interface {
	ExitWithContext(ctx context.Context, code int) error
}

// ContextualInitializer provides context-aware initialization.
type ContextualInitializer interface {
	InitWithContext(ctx context.Context) error
}

// ModernCommander combines all the modern interfaces for a full-featured command.
type ModernCommander interface {
	Namer
	Aliaser
	Describer
	Logger
	ContextualCommander
	ContextualExiter
	ContextualInitializer
}

// Ensure Cmd implements ModernCommander
var _ ModernCommander = (*Cmd)(nil)