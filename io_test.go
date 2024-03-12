package cleo

import (
	"fmt"
	"testing"

	"github.com/markbates/iox"
	"github.com/markbates/plugins"
	"github.com/stretchr/testify/require"
)

var _ plugins.IOSetable = &ioPlugin{}
var _ plugins.IOable = &ioPlugin{}

type ioPlugin struct {
	iox.IO
}

func (i *ioPlugin) SetStdio(oi iox.IO) error {
	if i == nil {
		return fmt.Errorf("nil ioPlugin")
	}

	i.IO = oi

	return nil
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

	var cmd *Cmd

	oi := iox.Discard()

	r.Error(cmd.SetStdio(oi))

	act := cmd.Stdio()
	r.Equal(iox.IO{}, act)

	cmd = &Cmd{
		Name: "main",
		Feeder: func() plugins.Plugins {
			return plugins.Plugins{
				stringPlug("mystring"),
			}
		},
	}

	r.NotEqual(oi, cmd.IO)

	cmd.SetStdio(oi)

	act = cmd.Stdio()
	r.Equal(oi, act)
}
