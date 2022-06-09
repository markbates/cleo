package cleo

import "fmt"

const (
	ErrNoCommand  = stringErr("no command specified")
	ErrNoCommands = stringErr("no commands registered")
)

type ErrUnknownCommand string

func (e ErrUnknownCommand) Error() string {
	return fmt.Sprintf("unknown command %q", string(e))
}

type stringErr string

func (s stringErr) Error() string {
	return string(s)
}
