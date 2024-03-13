package cleo

import (
	"fmt"
	"io/fs"
)

func (cmd *Cmd) FileSystem() (fs.FS, error) {
	if cmd == nil {
		return nil, fmt.Errorf("nil command")
	}

	cmd.mu.RLock()
	defer cmd.mu.RUnlock()

	if cmd.FS == nil {
		return nil, fmt.Errorf("fs.FS is nil")
	}

	return cmd.FS, nil
}

func (cmd *Cmd) SetFileSystem(cab fs.FS) error {
	if cmd == nil {
		return fmt.Errorf("nil command")
	}

	if cab == nil {
		return fmt.Errorf("fs.FS is nil")
	}

	cmd.mu.Lock()
	cmd.FS = cab
	cmd.mu.Unlock()

	return nil
}
