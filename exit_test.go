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
		Desc: "My Description",
		IO: iox.IO{
			Out: oi.Stdout(),
			In:  oi.Stdin(),
			Err: oi.Stderr(),
		},
		Commands: map[string]Commander{
			"abc": newEcho(t, "abc"),
			"xyz": newEcho(t, "xyz"),
		},
		Feeder: func() plugins.Plugins {
			return plugins.Plugins{
				String("mystring"),
			}
		},
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
  mystring               github.com/markbates/cleo.String

Error: boom`

	r.Equal(exp, act)
}
