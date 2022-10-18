package cleo

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/markbates/plugins"
	"github.com/stretchr/testify/require"
)

type fsPlugin struct {
	fs.FS
}

func (f *fsPlugin) SetFileSystem(cab fs.FS) {
	f.FS = cab
}

func (f fsPlugin) PluginName() string {
	return "fs"
}

func Test_Cmd_SetFileSystem(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	cab := fstest.MapFS{}
	fsp := &fsPlugin{}

	cmd := &Cmd{
		Name: "main",
		Feeder: func() plugins.Plugins {
			return plugins.Plugins{
				fsp,
			}
		},
	}

	r.Nil(cmd.FS)
	r.Nil(fsp.FS)

	cmd.SetFileSystem(cab)
	r.Equal(cab, cmd.FileSystem())
	r.Equal(cab, fsp.FS)

}
