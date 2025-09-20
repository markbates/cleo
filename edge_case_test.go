package cleo

import (
	"testing"

	"github.com/markbates/plugins"
	"github.com/stretchr/testify/require"
)

// Test ScopedPlugins edge cases
func TestCmd_ScopedPlugins_EdgeCases(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	// Test with nil feeder
	cmd := &Cmd{}
	plugs := cmd.ScopedPlugins()
	r.Nil(plugs)

	// Test with feeder returning empty plugins
	cmd = &Cmd{
		Feeder: func() plugins.Plugins {
			return plugins.Plugins{}
		},
	}
	plugs = cmd.ScopedPlugins()
	r.Nil(plugs)
}

// Test MarshalJSON error cases
func TestCmd_MarshalJSON_ErrorScenarios(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	// Test marshaling with complex data that might cause errors
	cmd := &Cmd{
		Name:    "test\x00\x01", // problematic characters
		Aliases: []string{"alias\x00", "alias\x01"},
	}

	// This should still work since Go's JSON marshaler handles these
	data, err := cmd.MarshalJSON()
	r.NoError(err)
	r.NotNil(data)
}

// Test WithSubCommand with existing command map
func TestWithSubCommand_ExistingMap(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	// Create command with existing Commands map
	cmd := &Cmd{
		Name: "parent",
		Commands: map[string]Commander{
			"existing": &testCommander{name: "existing"},
		},
	}

	// Add new sub-command
	opt := WithSubCommand("new", &testCommander{name: "new"})
	err := opt(cmd)
	r.NoError(err)

	// Should have both commands
	r.Len(cmd.Commands, 2)
	r.Contains(cmd.Commands, "existing")
	r.Contains(cmd.Commands, "new")
}