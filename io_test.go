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

	s := IO{}
	r.Equal(os.Stdin, s.Stdin())
	r.Equal(os.Stdout, s.Stdout())
	r.Equal(os.Stderr, s.Stderr())

	var out bytes.Buffer
	var err bytes.Buffer
	var in bytes.Reader

	s.In = &in
	r.Equal(&in, s.Stdin())

	s.Out = &out
	r.Equal(&out, s.Stdout())

	s.Err = &err
	r.Equal(&err, s.Stderr())
}
