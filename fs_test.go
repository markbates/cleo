package cleo

import (
	"io/fs"
	"testing"
	"testing/fstest"

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

	cmd := &Cmd{
		Name: "main",
	}

	r.Nil(cmd.FS)

	cmd.SetFileSystem(cab)
	r.Equal(cab, cmd.FileSystem())

}
