package cleo

import (
	"testing"
	"testing/fstest"

	"github.com/markbates/iox"
	"github.com/markbates/plugins"
	"github.com/stretchr/testify/require"
)

type neederPlug struct {
	plugins.Plugins
}

func (n neederPlug) PluginName() string {
	return "neederPlug"
}

func (n *neederPlug) WithPlugins(fn plugins.Feeder) {
	n.Plugins = fn()
}

type availPlug bool

func (a availPlug) PluginName() string {
	return "availPlug"
}

func (a availPlug) PluginAvailable(root string) bool {
	return bool(a)
}

func Test_Init(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	cab := fstest.MapFS{}
	oi := iox.Discard()

	iop := &ioPlugin{}
	fsp := &fsPlugin{}
	np := &neederPlug{}
	yes := availPlug(true)
	no := availPlug(false)

	c := &Cmd{}

	cmd := &Cmd{
		Name: "main",
		FS:   cab,
		IO:   oi,
		Plugins: plugins.Plugins{
			iop,
			fsp,
			np,
			yes,
			no,
			c,
		},
	}

	var i int
	err := Init(cmd, "foo", func(p plugins.Plugin) {
		i++
	})

	r.NoError(err)
	r.Equal(5, i)

	r.Equal(cab, cmd.FileSystem())
	r.Equal(cmd.FileSystem(), fsp.FS)

	r.Equal(oi, cmd.Stdio())
	r.Equal(oi, c.Stdio())
	r.Equal(oi, iop.IO)

	r.Equal(5, len(np.Plugins))
}
