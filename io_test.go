package cleo

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_StdIO(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	s := &StdIO{}
	r.Equal(os.Stdin, s.In())
	r.Equal(os.Stdout, s.Out())
	r.Equal(os.Stderr, s.Err())

	var out bytes.Buffer
	var err bytes.Buffer
	var in bytes.Reader

	s = WithIn(s, &in)
	r.Equal(&in, s.In())

	s = WithOut(s, &out)
	r.Equal(&out, s.Out())

	s = WithErr(s, &err)
	r.Equal(&err, s.Err())
}
