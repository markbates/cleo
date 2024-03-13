package cleo

import (
	"testing"

	"github.com/markbates/iox"
	"github.com/markbates/plugins"
	"github.com/markbates/plugins/plugtest"
	"github.com/stretchr/testify/require"
)

func Test_Cmd_IO(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	var cmd *Cmd

	oi := iox.Discard()

	r.Error(cmd.SetStdio(oi))

	act := cmd.Stdio()
	r.Equal(iox.IO{}, act)

	cmd = &Cmd{
		Name: "main",
		Feeder: func() plugins.Plugins {
			return plugins.Plugins{
				plugtest.StringPlugin("mystring"),
			}
		},
	}

	r.NotEqual(oi, cmd.IO)

	cmd.SetStdio(oi)

	act = cmd.Stdio()
	r.Equal(oi, act)
}
