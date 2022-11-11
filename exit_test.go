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

	fn := func() plugins.Plugins {
		return plugins.Plugins{
			newEcho(t, "abc"),
			newEcho(t, "xyz"),
		}
	}

	cmd := &Cmd{
		Name: "main",
		Desc: "My Description",
		IO: iox.IO{
			Out: oi.Stdout(),
			In:  oi.Stdin(),
			Err: oi.Stderr(),
		},
		Feeder: fn,
		ExitFn: func(i int) {
			r.Equal(code, i)
		},
	}

	app := &echoPlug{
		Cmd: cmd,
	}

	Exit(app, code, boom)

	act := buf.Err.String()
	act = strings.TrimSpace(act)

	// fmt.Println(act)

	exp := `$ main
------
*github.com/markbates/cleo.echoPlug

Available Commands:
  Command  Description
  -------  -----------
  abc      echo abc
  xyz      echo xyz

Using Plugins:
  Name      Description  Type
  ----      -----------  ----
  echo/abc  echo abc     *github.com/markbates/cleo.echoPlug
  echo/xyz  echo xyz     *github.com/markbates/cleo.echoPlug

Error: boom`

	r.Equal(exp, act)
}
