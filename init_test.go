package cleo

// type neederPlug struct {
// 	plugins.Plugins
// }

// func (n neederPlug) PluginName() string {
// 	return "neederPlug"
// }

// func (n *neederPlug) WithPlugins(fn plugins.FeederFn) {
// 	n.Plugins = fn()
// }

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
