package cleo

import (
	"sort"

	"github.com/markbates/iox"
	"github.com/markbates/plugins"
)

type AvailabilityChecker = plugins.AvailabilityChecker
type FSSetable = plugins.FSSetable
type IOSetable = iox.IOSetable
type Needer = plugins.Needer

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
func Init(cmd *Cmd, root string, fns ...func(p plugins.Plugin)) error {
	if cmd == nil {
		return ErrNoCommand
	}

	plugs := cmd.ScopedPlugins()

	cab := cmd.FileSystem()

	plugs = plugs.Available(root)

	sort.Sort(plugs)

	for _, p := range plugs {
		if ps, ok := p.(iox.IOSetable); ok {
			ps.SetStdio(cmd.IO)
		}

		if ps, ok := p.(FSSetable); ok {
			ps.SetFileSystem(cab)
		}

		if ps, ok := p.(Needer); ok {
			ps.WithPlugins(func() plugins.Plugins {
				return plugs
			})
		}

		for _, fn := range fns {
			fn(p)
		}
	}

	return nil
}
