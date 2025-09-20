package cleo

import (
	"context"
	"log/slog"
	"testing"
	"testing/fstest"

	"github.com/markbates/iox"
	"github.com/stretchr/testify/require"
)

func TestNewCmd(t *testing.T) {
	tests := []struct {
		name    string
		cmdName string
		opts    []Option
		wantErr bool
		errType error
	}{
		{
			name:    "empty name",
			cmdName: "",
			wantErr: true,
			errType: ErrEmptyCommandName,
		},
		{
			name:    "valid name only",
			cmdName: "test",
			wantErr: false,
		},
		{
			name:    "with stdio option",
			cmdName: "test",
			opts:    []Option{WithStdio(iox.Discard())},
			wantErr: false,
		},
		{
			name:    "with filesystem option",
			cmdName: "test",
			opts:    []Option{WithFS(fstest.MapFS{})},
			wantErr: false,
		},
		{
			name: "with multiple options",
			cmdName: "test",
			opts: []Option{
				WithStdio(iox.Discard()),
				WithFS(fstest.MapFS{}),
				WithDescription("test command"),
				WithAliases("t", "tst"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := require.New(t)

			cmd, err := NewCmd(tt.cmdName, tt.opts...)

			if tt.wantErr {
				r.Error(err)
				if tt.errType != nil {
					r.ErrorIs(err, tt.errType)
				}
				r.Nil(cmd)
			} else {
				r.NoError(err)
				r.NotNil(cmd)
				r.Equal(tt.cmdName, cmd.CmdName())
				r.NotNil(cmd.Commands)
				r.NotNil(cmd.Logger())
			}
		})
	}
}

func TestCmd_ExitWithContext(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *Cmd
		code    int
		wantErr bool
		errType error
	}{
		{
			name:    "nil command",
			setup:   func() *Cmd { return nil },
			code:    1,
			wantErr: true,
			errType: ErrNilCommand,
		},
		{
			name: "valid exit",
			setup: func() *Cmd {
				cmd, _ := NewCmd("test")
				return cmd
			},
			code:    0,
			wantErr: false,
		},
		{
			name: "exit with custom function",
			setup: func() *Cmd {
				cmd, _ := NewCmd("test", WithExitFunction(func(code int) error {
					if code != 0 {
						return ErrUnknownCommand("test error")
					}
					return nil
				}))
				return cmd
			},
			code:    1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := require.New(t)

			cmd := tt.setup()
			ctx := context.Background()
			err := cmd.ExitWithContext(ctx, tt.code)

			if tt.wantErr {
				r.Error(err)
				if tt.errType != nil {
					r.ErrorIs(err, tt.errType)
				}
			} else {
				r.NoError(err)
			}
		})
	}
}

func TestCmd_InitWithContext(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *Cmd
		wantErr bool
		errType error
	}{
		{
			name:    "nil command",
			setup:   func() *Cmd { return nil },
			wantErr: true,
			errType: ErrNilCommand,
		},
		{
			name: "valid init",
			setup: func() *Cmd {
				cmd, _ := NewCmd("test", 
					WithStdio(iox.Discard()),
					WithFS(fstest.MapFS{}),
				)
				return cmd
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := require.New(t)

			cmd := tt.setup()
			ctx := context.Background()
			err := cmd.InitWithContext(ctx)

			if tt.wantErr {
				r.Error(err)
				if tt.errType != nil {
					r.ErrorIs(err, tt.errType)
				}
			} else {
				r.NoError(err)
			}
		})
	}
}

func TestCmd_Logger(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	// Test nil command
	var cmd *Cmd
	logger := cmd.Logger()
	r.Equal(slog.Default(), logger)

	// Test command without logger
	cmd, err := NewCmd("test")
	r.NoError(err)
	logger = cmd.Logger()
	r.NotNil(logger)

	// Test command with custom logger
	customLogger := slog.Default()
	cmd, err = NewCmd("test", WithLogger(customLogger))
	r.NoError(err)
	logger = cmd.Logger()
	r.Equal(customLogger, logger)
}