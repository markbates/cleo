package cleo

import (
	"github.com/markbates/iox"
)

func (cmd *Cmd) Stdio() iox.IO {
	if cmd == nil {
		return iox.IO{}
	}

	cmd.RLock()
	defer cmd.RUnlock()

	return cmd.IO
}

func (cmd *Cmd) SetStdio(oi iox.IO) {
	if cmd == nil {
		return
	}

	cmd.Lock()
	cmd.IO = oi
	cmd.Unlock()
}
