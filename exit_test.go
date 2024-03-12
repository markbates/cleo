package cleo

import (
	"fmt"
	"strings"
	"testing"

	"github.com/markbates/iox"
	"github.com/markbates/plugins"
	"github.com/stretchr/testify/require"
)

func Test_Exit(t *testing.T) {
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
				stringPlug("mystring"),
			}
		},
		ExitFn: func(i int) error {
			r.Equal(code, i)
			return nil
		},
	}

	// app := &echoPlug{}

	Exit(cmd, code, boom)

	act := buf.Err.String()
	act = strings.TrimSpace(act)

	// fmt.Println(act)

	exp := `$ main
------
*github.com/markbates/cleo.Cmd

Available Commands:
  Command  Description
  -------  -----------
  abc      echo abc
  xyz      echo xyz

Using Plugins:
  Name      Description  Type
  ----      -----------  ----
  mystring               github.com/markbates/cleo.stringPlug

Error: boom`

	r.Equal(exp, act)
}

func Test_Cmd_Exit(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	var cmd *Cmd

	err := cmd.Exit(42)
	r.Error(err)

	ep := &exiterPlug{}
	plugs := plugins.Plugins{
		ep,
	}

	cmd = &Cmd{
		Feeder: func() plugins.Plugins {
			return plugs
		},
	}
	err = cmd.Exit(42)
	r.NoError(err)

	r.Equal(42, ep.Code)

	err = cmd.Exit(-1)
	r.Error(err)
}
