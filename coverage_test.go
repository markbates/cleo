package cleo

import (
	"context"
	"io"
	"log/slog"
	"os"
	"testing"
	"testing/fstest"

	"github.com/markbates/iox"
	"github.com/markbates/plugins"
	"github.com/stretchr/testify/require"
)

// Test Env.Set method (currently 0% coverage)
func TestEnv_Set(t *testing.T) {
	tests := []struct {
		name    string
		env     *Env
		key     string
		value   string
		wantErr bool
	}{
		{
			name:  "nil env sets to os env",
			env:   nil,
			key:   "TEST_NIL_ENV",
			value: "test_value",
		},
		{
			name:  "empty env sets to internal map",
			env:   &Env{},
			key:   "TEST_EMPTY_ENV",
			value: "test_value",
		},
		{
			name: "existing env sets to internal map",
			env: &Env{
				data: map[string]string{
					"EXISTING": "existing_value",
				},
			},
			key:   "NEW_KEY",
			value: "new_value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			// Save original env value if testing with nil env
			var originalValue string
			var wasSet bool
			if tt.env == nil {
				originalValue, wasSet = os.LookupEnv(tt.key)
				defer func() {
					if wasSet {
						os.Setenv(tt.key, originalValue)
					} else {
						os.Unsetenv(tt.key)
					}
				}()
			}

			tt.env.Set(tt.key, tt.value)

			if tt.env == nil {
				// Check OS environment
				actual := os.Getenv(tt.key)
				r.Equal(tt.value, actual)
			} else {
				// Check internal map
				r.NotNil(tt.env.data)
				r.Equal(tt.value, tt.env.data[tt.key])
			}
		})
	}
}

// Test ErrUnknownCommand.Error method (currently 0% coverage)
func TestErrUnknownCommand_Error(t *testing.T) {
	tests := []struct {
		name     string
		cmdName  string
		expected string
	}{
		{
			name:     "simple command name",
			cmdName:  "test",
			expected: `unknown command "test"`,
		},
		{
			name:     "command with spaces",
			cmdName:  "test command",
			expected: `unknown command "test command"`,
		},
		{
			name:     "empty command name",
			cmdName:  "",
			expected: `unknown command ""`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			err := ErrUnknownCommand(tt.cmdName)
			r.Equal(tt.expected, err.Error())
		})
	}
}

// Test WithPluginFeeder option (currently 0% coverage)
func TestWithPluginFeeder(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	feeder := func() plugins.Plugins {
		return plugins.Plugins{
			&testPlugin{name: "test-plugin"},
		}
	}

	cmd, err := NewCmd("test", WithPluginFeeder(feeder))
	r.NoError(err)
	r.NotNil(cmd)

	plugs := cmd.ScopedPlugins()
	r.Len(plugs, 1)

	// Test nil command
	var nilCmd *Cmd
	opt := WithPluginFeeder(feeder)
	err = opt(nilCmd)
	r.ErrorIs(err, ErrNilCommand)
}

// Test WithSubCommand option (currently 0% coverage)
func TestWithSubCommand(t *testing.T) {
	tests := []struct {
		name        string
		cmdName     string
		subCmdName  string
		subCmd      Commander
		wantErr     bool
		expectedErr error
	}{
		{
			name:       "valid sub-command",
			cmdName:    "parent",
			subCmdName: "child",
			subCmd:     &testCommander{name: "child"},
			wantErr:    false,
		},
		{
			name:        "empty sub-command name",
			cmdName:     "parent",
			subCmdName:  "",
			subCmd:      &testCommander{name: "child"},
			wantErr:     true,
			expectedErr: ErrEmptyCommandName,
		},
		{
			name:        "nil sub-command",
			cmdName:     "parent",
			subCmdName:  "child",
			subCmd:      nil,
			wantErr:     true,
			expectedErr: ErrNilSubCommand,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := require.New(t)

			cmd, err := NewCmd(tt.cmdName, WithSubCommand(tt.subCmdName, tt.subCmd))

			if tt.wantErr {
				r.Error(err)
				if tt.expectedErr != nil {
					r.ErrorIs(err, tt.expectedErr)
				}
				r.Nil(cmd)
			} else {
				r.NoError(err)
				r.NotNil(cmd)
				r.Contains(cmd.Commands, tt.subCmdName)
				r.Equal(tt.subCmd, cmd.Commands[tt.subCmdName])
			}
		})
	}

	// Test nil command
	t.Run("nil command", func(t *testing.T) {
		r := require.New(t)
		var nilCmd *Cmd
		opt := WithSubCommand("test", &testCommander{name: "test"})
		err := opt(nilCmd)
		r.ErrorIs(err, ErrNilCommand)
	})
}

// Test option error cases more thoroughly
func TestOptionsErrorCases(t *testing.T) {
	tests := []struct {
		name        string
		option      Option
		expectedErr error
	}{
		{
			name:        "WithStdio with nil cmd",
			option:      WithStdio(iox.Discard()),
			expectedErr: ErrNilCommand,
		},
		{
			name:        "WithFS with nil cmd",
			option:      WithFS(fstest.MapFS{}),
			expectedErr: ErrNilCommand,
		},
		{
			name:        "WithAliases with nil cmd",
			option:      WithAliases("a", "b"),
			expectedErr: ErrNilCommand,
		},
		{
			name:        "WithDescription with nil cmd",
			option:      WithDescription("test"),
			expectedErr: ErrNilCommand,
		},
		{
			name:        "WithExitFunction with nil cmd",
			option:      WithExitFunction(func(int) error { return nil }),
			expectedErr: ErrNilCommand,
		},
		{
			name:        "WithLogger with nil cmd",
			option:      WithLogger(slog.Default()),
			expectedErr: ErrNilCommand,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := require.New(t)

			var nilCmd *Cmd
			err := tt.option(nilCmd)
			r.ErrorIs(err, tt.expectedErr)
		})
	}
}

// Test WithFS with nil filesystem separately
func TestWithFS_NilFilesystem(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	cmd, err := NewCmd("test", WithFS(nil))
	r.ErrorIs(err, ErrNilFS)
	r.Nil(cmd)
}

// Test more comprehensive NewCmd scenarios
func TestNewCmd_ErrorCases(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	// Test with failing option
	failingOption := func(cmd *Cmd) error {
		return ErrNilFS
	}

	cmd, err := NewCmd("test", failingOption)
	r.Error(err)
	r.Nil(cmd)
	r.Contains(err.Error(), "applying option")
}

// Test MarshalJSON error cases
func TestCmd_MarshalJSON_ErrorCases(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	// Test nil command
	var nilCmd *Cmd
	data, err := nilCmd.MarshalJSON()
	r.ErrorIs(err, ErrNilCommand)
	r.Nil(data)
}

// Test Exit function more thoroughly
func TestExit_MoreCases(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	// Test with nil error (should return early)
	mockCmd := &mockStdioer{}
	Exit(mockCmd, 1, nil)
	r.False(mockCmd.printCalled, "Print should not be called when err is nil")

	// Test with code -1 (should print usage only)
	mockCmd = &mockStdioer{}
	Exit(mockCmd, -1, ErrNoCommand)
	r.True(mockCmd.printCalled, "Print should be called with code -1")
	r.False(mockCmd.exitCalled, "Exit should not be called with code -1")

	// Test with normal error code
	mockCmd = &mockStdioer{}
	Exit(mockCmd, 1, ErrNoCommand)
	r.True(mockCmd.printCalled, "Print should be called with normal error")
	r.True(mockCmd.exitCalled, "Exit should be called with normal error")
}

// Helper types for testing

type testPlugin struct {
	name string
}

func (p *testPlugin) PluginName() string {
	return p.name
}

type testCommander struct {
	name string
}

func (c *testCommander) CmdName() string {
	return c.name
}

func (c *testCommander) Main(ctx context.Context, pwd string, args []string) error {
	return nil
}

func (c *testCommander) PluginName() string {
	return c.name
}

type mockStdioer struct {
	printCalled bool
	exitCalled  bool
}

func (m *mockStdioer) Stderr() io.Writer {
	m.printCalled = true
	return iox.Discard().Stderr()
}

func (m *mockStdioer) Stdout() io.Writer {
	return iox.Discard().Stdout()
}

func (m *mockStdioer) Stdin() io.Reader {
	return iox.Discard().Stdin()
}

func (m *mockStdioer) Exit(code int) error {
	m.exitCalled = true
	return nil
}

func (m *mockStdioer) PluginName() string {
	return "mockStdioer"
}