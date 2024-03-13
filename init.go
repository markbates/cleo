package cleo

import (
	"fmt"

	"github.com/markbates/plugins"
)

// Init will initialize the command.
// It should be called before the command is used.
func (cmd *Cmd) Init() error {
	if cmd == nil {
		return fmt.Errorf("nil command")
	}

	cmd.mu.Lock()
	cab := cmd.FS
	cmd.mu.Unlock()

	plugFn := cmd.PluginFeeder()
	plugs := plugFn()

	// plugins.FSSetable
	fss := plugins.ByType[plugins.FSSetable](plugs)
	for _, fs := range fss {
		if err := fs.SetFileSystem(cab); err != nil {
			return err
		}
	}

	// plugins.Needer
	needs := plugins.ByType[plugins.Needer](plugs)
	for _, n := range needs {
		if err := n.WithPlugins(plugFn); err != nil {
			return err
		}
	}

	// plugins.IOSetable
	ios := plugins.ByType[plugins.IOSetable](plugs)
	for _, io := range ios {
		if err := io.SetStdio(cmd.IO); err != nil {
			return err
		}
	}

	return nil
}
