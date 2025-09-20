package cleo

import (
	"errors"
	"fmt"
)

// Common errors
var (
	ErrNoCommand        = errors.New("no command specified")
	ErrNoCommands       = errors.New("no commands registered")
	ErrNilCommand       = errors.New("nil command")
	ErrNilFS            = errors.New("fs.FS is nil")
	ErrEmptyCommandName = errors.New("empty command name")
	ErrNilSubCommand    = errors.New("nil sub-command")
)

type ErrUnknownCommand string

func (e ErrUnknownCommand) Error() string {
	return fmt.Sprintf("unknown command %q", string(e))
}
