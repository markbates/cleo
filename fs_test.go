package cleo

import (
	"fmt"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/markbates/plugins"
	"github.com/stretchr/testify/require"
)

var _ plugins.FSSetable = &fsPlugin{}
var _ plugins.FSable = &fsPlugin{}

type fsPlugin struct {
	fs.FS
}

func (f *fsPlugin) FileSystem() (fs.FS, error) {
	if f == nil {
		return nil, fmt.Errorf("fsPlugin is nil")
	}

	if f.FS == nil {
		return nil, fmt.Errorf("fs.FS is nil")
	}

	return f.FS, nil
}

func (f *fsPlugin) SetFileSystem(cab fs.FS) error {
	if f == nil {
		return fmt.Errorf("fsPlugin is nil")
	}

	if cab == nil {
		return fmt.Errorf("fs.FS is nil")
	}

	f.FS = cab

	return nil
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

	err := cmd.SetFileSystem(cab)
	r.NoError(err)

	kab, err := cmd.FileSystem()
	r.NoError(err)

	r.Equal(cab, kab)

}
