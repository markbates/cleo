package cleo

import (
	"testing"

	"github.com/markbates/iox"
	"github.com/markbates/plugins"
	"github.com/stretchr/testify/require"
)

type ioPlugin struct {
	iox.IO
}

func (i *ioPlugin) SetStdio(oi iox.IO) {
	i.IO = oi
}

func (i ioPlugin) Stdio() iox.IO {
	return i.IO
}

func (i ioPlugin) PluginName() string {
	return "ioPlugin"
}

func Test_Cmd_IO(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	oi := iox.Discard()

	cmd := &Cmd{
		Name: "main",
		Feeder: func() plugins.Plugins {
			return plugins.Plugins{
				String("mystring"),
			}
		},
	}

	r.NotEqual(oi, cmd.IO)

	cmd.SetStdio(oi)
	r.Equal(oi, cmd.IO)
}
