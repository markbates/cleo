package cleo

import (
	"context"
	"strings"
	"testing"

	"github.com/markbates/plugins"
	"github.com/markbates/plugins/plugtest"
	"github.com/stretchr/testify/require"
)

func Test_Cmd_ScopedPlugins(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	var cmd *Cmd
	plugs := cmd.ScopedPlugins()
	r.Len(plugs, 0)

	cmd = &Cmd{}
	plugs = cmd.ScopedPlugins()
	r.Len(plugs, 0)

	cmd.Feeder = func() plugins.Plugins {
		return plugins.Plugins{
			plugtest.StringPlugin("mystring"),
		}
	}

	plugs = cmd.ScopedPlugins()
	r.Len(plugs, 1)

}

func Test_Cmd_SubCommands(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	var cmd *Cmd
	cmds := cmd.SubCommands()
	r.Len(cmds, 0)

	fn := func() plugins.Plugins {
		return plugins.Plugins{
			plugtest.StringPlugin("mystring"),
		}
	}

	cmd = &Cmd{
		Name: "main",
		Commands: map[string]Commander{
			"abc": newCleoPlug(t, "abc"),
			"xyz": newCleoPlug(t, "xyz"),
		},
		Feeder: fn,
	}

	cmds = cmd.SubCommands()
	r.Len(cmds, 2)

	c, ok := cmds[0].(*cleoPlug)
	r.True(ok)
	r.Equal("abc", c.Name)

	c, ok = cmds[1].(*cleoPlug)
	r.True(ok)
	r.Equal("xyz", c.Name)
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

	var cmd *Cmd

	err := cmd.Init()
	r.Error(err)

	cmd = &Cmd{}

	err = cmd.Init()
	r.NoError(err)
}

func Test_Cmd_CmdName(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	var cmd *Cmd
	r.Equal("", cmd.CmdName())

	cmd = &Cmd{Name: "main"}
	r.Equal("main", cmd.CmdName())
}

func Test_Cmd_CmdAliases(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	var cmd *Cmd
	r.Nil(cmd.CmdAliases())

	cmd = &Cmd{Aliases: []string{"a", "b"}}
	r.Equal([]string{"a", "b"}, cmd.CmdAliases())
}

func Test_Cmd_String(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	var cmd *Cmd
	r.Equal("", cmd.String())

	cmd = &Cmd{Name: "main"}

	act := cmd.String()

	exp := `{"aliases":null,"name":"main","plugins":null,"stdio":{}}`
	r.Equal(exp, act)
}

func Test_Cmd_MarshalJSON(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	var cmd *Cmd
	_, err := cmd.MarshalJSON()
	r.Error(err)

	cmd = &Cmd{
		Name:    "main",
		Aliases: []string{"a", "b"},
		Feeder: func() plugins.Plugins {
			return plugins.Plugins{
				plugtest.StringPlugin("mystring"),
			}
		},
		Desc: "My Description",
	}

	b, err := cmd.MarshalJSON()
	r.NoError(err)

	act := string(b)
	act = strings.TrimSpace(act)

	// fmt.Println(act)

	exp := `{"aliases":["a","b"],"name":"main","plugins":["mystring"],"stdio":{}}`

	r.Equal(exp, act)

}

func Test_Cmd_Main(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	var cmd *Cmd

	ctx := context.Background()

	err := cmd.Main(ctx, "", []string{})
	r.Error(err)

	cmd = &Cmd{}

	err = cmd.Main(ctx, "", []string{})
	r.Error(err)

}

func Test_Cmd_Description(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	var cmd *Cmd

	r.Equal("", cmd.Description())

	exp := "My Description"
	cmd = &Cmd{
		Desc: exp,
	}

	r.Equal(exp, cmd.Description())
}

func Test_Cmd_PluginFeeder(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	var cmd *Cmd

	fn := cmd.PluginFeeder()
	r.NotNil(fn)

	plugs := fn()
	r.Len(plugs, 0)

	cmd = &Cmd{}

	fn = cmd.PluginFeeder()
	r.NotNil(fn)

	plugs = fn()
	r.Len(plugs, 0)

	cmd.Feeder = func() plugins.Plugins {
		return plugins.Plugins{
			plugtest.StringPlugin("mystring"),
		}
	}

	fn = cmd.PluginFeeder()
	r.NotNil(fn)

	plugs = fn()
	r.Len(plugs, 1)

	cmd.Feeder = func() plugins.Plugins {
		return plugins.Plugins{
			plugtest.StringPlugin("mystring"),
			&cleoPlug{
				Plugins: func() plugins.Plugins {
					return plugins.Plugins{
						plugtest.StringPlugin("another string"),
					}
				},
			},
		}
	}

	fn = cmd.PluginFeeder()
	r.NotNil(fn)

	plugs = fn()
	r.Len(plugs, 3)

}
