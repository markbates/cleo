package cleo

import "io/fs"

func (cmd *Cmd) FileSystem() fs.FS {
	if cmd == nil {
		return nil
	}

	cmd.RLock()
	defer cmd.RUnlock()
	return cmd.FS
}

func (cmd *Cmd) SetFileSystem(cab fs.FS) {
	if cmd == nil {
		return
	}

	plugs := cmd.ScopedPlugins()

	cmd.Lock()
	defer cmd.Unlock()

	cmd.FS = cab

	plugs.SetFileSystem(cab)
}
