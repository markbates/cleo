package cleo

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Switcher_Switch(t *testing.T) {
	t.Parallel()

	var logFn SwitchFn
	logFn = func(rt *Runtime) error {
		fmt.Fprintf(rt.Stdout, "%s\n", rt)

		if next, ok := rt.Next(); ok {
			return logFn(next)
		}

		return nil
	}

	defSw := NewSwitcher()
	defSw.SetDefault(logFn)

	simpleSw := NewSwitcher()
	simpleSw.Set("hello", logFn)

	recurseSw := NewSwitcher()
	recurseSw.Set("hello", logFn)
	recurseSw.Set("world", logFn)

	table := []struct {
		args []string
		err  bool
		exp  string
		name string
		sw   *Switcher
	}{
		{name: "default switch", sw: defSw, exp: "Runtime: test\n"},
		{name: "no args, no def", sw: NewSwitcher(), err: true},
		{name: "recurse switch", sw: recurseSw, args: []string{"hello", "world"}, exp: "Runtime: hello [world]\nRuntime: world\n"},
		{name: "simple switch", sw: simpleSw, args: []string{"hello"}, exp: "Runtime: hello\n"},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			rt, err := NewRuntime("test", tt.args)
			r.NoError(err)

			bb := &bytes.Buffer{}
			rt.Stdout = bb

			err = tt.sw.Switch(rt)
			if tt.err {
				r.Error(err)
				return
			}

			r.NoError(err)

			act := bb.String()
			exp := tt.exp

			r.Equal(exp, act)

		})
	}
}
