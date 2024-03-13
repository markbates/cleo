package cleo

import (
	"testing"
	"testing/fstest"

	"github.com/markbates/iox"
	"github.com/markbates/plugins"
	"github.com/markbates/plugins/plugtest"
	"github.com/stretchr/testify/require"
)

func Test_CMD_Init(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	var cmd *Cmd
	r.Error(cmd.Init())

	cab := fstest.MapFS{}

	cmd = &Cmd{}

	r.NoError(cmd.Init())

	t.Run("FSSetable", func(st *testing.T) {
		r := require.New(st)

		fsp := &plugtest.FSable{}

		cmd := &Cmd{
			FS: cab,
			Feeder: func() plugins.Plugins {
				return plugins.Plugins{
					fsp,
				}
			},
		}

		r.NoError(cmd.Init())

		kab, err := fsp.FileSystem()
		r.NoError(err)

		r.Equal(cab, kab)

		cmd.FS = nil
		r.Error(cmd.Init())
	})

	t.Run("Needer", func(st *testing.T) {
		r := require.New(st)

		np := &plugtest.Needer{}

		cmd := &Cmd{
			Feeder: func() plugins.Plugins {
				return plugins.Plugins{
					np,
				}
			},
		}

		r.NoError(cmd.Init())

		r.NotNil(np.Fn)

		np = nil

		cmd.Feeder = func() plugins.Plugins {
			return plugins.Plugins{
				np,
			}
		}

		r.Error(cmd.Init())
	})

	t.Run("IOSetable", func(st *testing.T) {
		r := require.New(st)

		oi := iox.Discard()

		iop := &plugtest.IO{}

		cmd := &Cmd{
			IO: oi,
			Feeder: func() plugins.Plugins {
				return plugins.Plugins{
					iop,
				}
			},
		}

		r.NoError(cmd.Init())

		r.Equal(oi, iop.IO)

		iop = nil

		cmd.Feeder = func() plugins.Plugins {
			return plugins.Plugins{
				iop,
			}
		}

		r.Error(cmd.Init())
	})
}
