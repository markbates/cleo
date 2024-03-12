package cleo

import (
	"fmt"

	"github.com/markbates/iox"
)

func (cmd *Cmd) Stdio() iox.IO {
	if cmd == nil {
		return iox.IO{}
	}

	cmd.mu.RLock()
	defer cmd.mu.RUnlock()

	return cmd.IO
}

func (cmd *Cmd) SetStdio(oi iox.IO) error {
	if cmd == nil {
		return fmt.Errorf("nil command")
	}

	cmd.mu.Lock()
	cmd.IO = oi
	cmd.mu.Unlock()

	return nil
}
