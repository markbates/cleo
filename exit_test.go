package cleo

import (
	"fmt"
	"strings"
	"testing"

	"github.com/markbates/iox"
	"github.com/markbates/plugins"
	"github.com/stretchr/testify/require"
)

func Test_Cmd_Exit(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	boom := fmt.Errorf("boom")
	code := 42

	buf := &iox.Buffer{}
	oi := buf.IO()

	cmd := &Cmd{
		Name: "main",
		IO: iox.IO{
			Out: oi.Stdout(),
			In:  oi.Stdin(),
			Err: oi.Stderr(),
		},
		Plugins: plugins.Plugins{
			newEcho(t, "abc"),
			newEcho(t, "xyz"),
		},
		ExitFn: func(i int) {
			r.Equal(code, i)
		},
	}

	cmd.Exit(code, boom)

	act := buf.Err.String()
	act = strings.TrimSpace(act)

	// fmt.Println(act)

	exp := `$ main
------
*github.com/markbates/cleo.Cmd

Available Commands:
  Command  Description
  -------  -----------
  abc
  xyz

Using Plugins:
  Name              Description  Type
  ----              -----------  ----
  *cleo.Cmd (main)               *github.com/markbates/cleo.Cmd
  echo                           *github.com/markbates/cleo.echoPlug

Error: boom`

	r.Equal(exp, act)
}
