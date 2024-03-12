package cleo

import (
	"fmt"
	"testing"
	"testing/fstest"

	"github.com/markbates/iox"
	"github.com/markbates/plugins"
	"github.com/stretchr/testify/require"
)

var _ plugins.Needer = &neederPlug{}

type neederPlug struct {
	Plugins plugins.Plugins
}

func (n neederPlug) PluginName() string {
	return "neederPlug"
}

func (n *neederPlug) WithPlugins(fn plugins.FeederFn) error {
	if n == nil {
		return fmt.Errorf("nil neederPlug")
	}

	if fn == nil {
		return fmt.Errorf("nil FeederFn")
	}

	n.Plugins = fn()
	return nil
}

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

		fsp := &fsPlugin{}

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

		np := &neederPlug{}

		cmd := &Cmd{
			Feeder: func() plugins.Plugins {
				return plugins.Plugins{
					np,
				}
			},
		}

		r.NoError(cmd.Init())

		r.NotNil(np.Plugins)

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

		iop := &ioPlugin{}

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

// type availPlug bool

// func (a availPlug) PluginName() string {
// 	return "availPlug"
// }

// func (a availPlug) PluginAvailable(root string) bool {
// 	return bool(a)
// }

// func Test_Init(t *testing.T) {
// 	t.Parallel()
// 	r := require.New(t)

// 	cab := fstest.MapFS{}
// 	oi := iox.Discard()

// 	iop := &ioPlugin{}
// 	fsp := &fsPlugin{}
// 	np := &neederPlug{}
// 	yes := availPlug(true)
// 	no := availPlug(false)

// 	cmd := &Cmd{
// 		Name: "main",
// 		FS:   cab,
// 		IO:   oi,
// 	}

// 	fn := func() plugins.Plugins {
// 		return plugins.Plugins{
// 			iop,
// 			fsp,
// 			np,
// 			yes,
// 			no,
// 		}
// 	}

// 	cmd.Feeder = fn

// 	var i int
// 	err := Init(cmd, "foo", func(p plugins.Plugin) {
// 		i++
// 	})

// 	r.NoError(err)
// 	r.Equal(4, i)

// 	kab, err := cmd.FileSystem()
// 	r.NoError(err)

// 	r.Equal(cab, kab)

// 	r.Equal(oi, cmd.Stdio())
// 	r.Equal(oi, iop.IO)

// 	r.Equal(5, len(np.Plugins))
// }
