package cleo

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func Test_Cmd_SetFileSystem(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	cab := fstest.MapFS{}

	var cmd *Cmd
	r.Error(cmd.SetFileSystem(cab))

	cmd = &Cmd{
		Name: "main",
	}

	r.Nil(cmd.FS)
	r.Error(cmd.SetFileSystem(nil))

	err := cmd.SetFileSystem(cab)
	r.NoError(err)

	kab, err := cmd.FileSystem()
	r.NoError(err)

	r.Equal(cab, kab)

}

func Test_Cmd_FileSystem(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	var cmd *Cmd
	_, err := cmd.FileSystem()
	r.Error(err)

	cmd = &Cmd{
		Name: "main",
	}

	_, err = cmd.FileSystem()
	r.Error(err)

	cab := fstest.MapFS{}
	cmd.FS = cab

	kab, err := cmd.FileSystem()
	r.NoError(err)

	r.Equal(cab, kab)
}
