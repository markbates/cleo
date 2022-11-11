package cleo

import (
	"context"
	"os"
	"os/signal"
	"testing"
	"time"

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

func Test_Cmd_Init(t *testing.T) {
	t.Parallel()
	r := require.New(t)
	cmd := newEcho(t, "main")
	cmd.Feeder = func() plugins.Plugins {
		return plugins.Plugins{
			newEcho(t, "abc"),
			newEcho(t, "xyz"),
			String("mystring"),
		}
	}

	plugs := cmd.Plugins()
	plugs = append(plugs, newEcho(t, "123"), cmd)

	// your plugins here:
	// cmd.Plugins = append(cmd.Plugins, ...)
	fn := func() plugins.Plugins {
		return plugs
	}

	cmd.Feeder = fn

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	ctx, cancel = signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	r.NoError(Init(cmd.Cmd, "."))

	err := cmd.Main(ctx, ".", []string{"main", "abc", "hello"})
	r.NoError(err)

}
