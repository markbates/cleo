package cleo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewRuntime(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	in := []string{"json", "-i", "foo.json"}
	rt, err := NewRuntime("test", in)
	r.NoError(err)
	r.NotNil(rt)

	r.NotNil(rt.Cab)
	r.NotNil(rt.FlagSet)
	r.NotNil(rt.Stderr)
	r.NotNil(rt.Stdin)
	r.NotNil(rt.Stdout)

	r.Equal(in, rt.Args)
}

func Test_Runtime_Next(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	in := []string{"a", "b", "c"}

	a, err := NewRuntime("test", in)
	r.NoError(err)
	r.NotNil(a)

	r.Equal(a.Args, in)

	b, ok := a.Next()
	r.True(ok)
	r.Equal(b.Args, in[1:])

	c, ok := b.Next()
	r.True(ok)
	r.Equal(c.Args, in[2:])

	d, ok := c.Next()
	r.True(ok)
	r.Empty(d.Args)

	_, ok = d.Next()
	r.False(ok)
}
