package cleo

import (
	"testing"

	"github.com/markbates/plugins"
	"github.com/stretchr/testify/require"
)

func Test_Cmd_ScopedPlugins(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	cmd := &Cmd{
		Name: "main",
	}

	fn := func() plugins.Plugins {
		return plugins.Plugins{
			newEcho(t, "abc"),
			newEcho(t, "xyz"),
			String("mystring"),
			cmd,
		}
	}

	cmd.Feeder = fn

	scoped := cmd.ScopedPlugins()
	r.Len(scoped, 3)

}

func Test_Cmd_SubCommands(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	fn := func() plugins.Plugins {
		return plugins.Plugins{
			newEcho(t, "abc"),
			newEcho(t, "xyz"),
			String("mystring"),
		}
	}

	cmd := &Cmd{
		Name:   "main",
		Feeder: fn,
	}

	cmds := cmd.SubCommands()
	r.Len(cmds, 2)

	c, ok := cmds[0].(*echoPlug)
	r.True(ok)
	r.Equal("abc", c.CmdName())

	c, ok = cmds[1].(*echoPlug)
	r.True(ok)
	r.Equal("xyz", c.CmdName())
}

func Test_Cmd_PluginName(t *testing.T) {
	t.Parallel()

	table := []struct {
		cmd *Cmd
		exp string
	}{
		{&Cmd{Name: "main"}, "*cleo.Cmd (main)"},
		{nil, "*cleo.Cmd (?)"},
	}

	for _, tc := range table {
		t.Run(tc.exp, func(st *testing.T) {
			r := require.New(st)

			r.Equal(tc.exp, tc.cmd.PluginName())
		})
	}

}
