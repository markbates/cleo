package cleo

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Command_Main(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	cmd := &Cmd{}

	ctx := context.Background()
	err := cmd.Main(ctx, ".", nil)
	r.Equal(ErrNoCommands, err)

	cmd.Add("foo", newEcho(t))

	err = cmd.Main(ctx, ".", nil)
	r.Equal(ErrNoCommand, err)

	err = cmd.Main(ctx, ".", []string{"bar"})
	r.Equal(ErrUnknownCommand("bar"), err)

	bb := &bytes.Buffer{}
	cmd.IO.Out = bb

	err = cmd.Main(ctx, ".", []string{"foo", "1", "2", "3"})
	r.NoError(err)

	act := bb.String()
	exp := "[1 2 3]"

	r.Equal(exp, act)
}
