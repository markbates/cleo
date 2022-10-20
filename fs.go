package cleo

import (
	"io/fs"
)

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

	cmd.Lock()
	cmd.FS = cab
	cmd.Unlock()
}
