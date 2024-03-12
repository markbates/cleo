package cleo

// Init is a helper function that will
// initialize a Cmd with the plugins
// that are available for the given root
// directory.
// The following are the interfaces that
// are called on the ScopedPlugins of the
// given Cmd:
//
//	AvalabilityChecker
//	FSSetable
//	IOSetable
//	Needer
// func Init(cmd *Cmd, root string, fns ...func(p plugins.Plugin)) error {
// 	if cmd == nil {
// 		return ErrNoCommand
// 	}

// 	plugs := cmd.ScopedPlugins()

// 	cab, err := cmd.FileSystem()
// 	if err != nil {
// 		return err
// 	}

// 	plugs = plugs.Available(root)

// 	sort.Sort(plugs)

// 	for _, p := range plugs {
// 		if ps, ok := p.(iox.IOSetable); ok {
// 			if err := ps.SetStdio(cmd.IO); err != nil {
// 				return err
// 			}
// 		}

// 		if ps, ok := p.(FSSetable); ok {
// 			if err := ps.SetFileSystem(cab); err != nil {
// 				return err
// 			}
// 		}

// 		if ps, ok := p.(Needer); ok {
// 			err := ps.WithPlugins(func() plugins.Plugins {
// 				return plugs
// 			})
// 			if err != nil {
// 				return err
// 			}
// 		}

// 		for _, fn := range fns {
// 			fn(p)
// 		}
// 	}

// 	return nil
// }
