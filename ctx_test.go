package cleo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_NewContext(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	outer, cancel := context.WithTimeout(context.Background(), 2*time.Millisecond)
	defer cancel()

	inner, cancel := NewContext(outer)
	defer cancel()

	<-inner.Done()

	err := inner.Err()
	r.Equal(context.DeadlineExceeded, err)
}

func Test_ContextWithTimeout(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	outer, cancel := context.WithCancel(context.Background())
	defer cancel()

	inner, cancel := ContextWithTimeout(outer, 2*time.Millisecond)
	defer cancel()

	<-inner.Done()

	err := inner.Err()
	r.Equal(context.DeadlineExceeded, err)
}

func Test_Context_Error(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	outer, cancel := context.WithTimeout(context.Background(), 2*time.Millisecond)
	defer cancel()

	inner, cancel := NewContext(outer)
	defer cancel()

	<-inner.Done()

	err := inner.Err()
	r.Error(err)
	r.Equal(context.DeadlineExceeded, err)
}

func Test_Context_SetErr(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	boom := fmt.Errorf("boom")

	outer, cancel := context.WithCancel(context.Background())
	defer cancel()

	inner, cancel := ContextWithTimeout(outer, 2*time.Millisecond)
	defer cancel()

	inner.SetErr(boom)

	<-inner.Done()

	err := inner.Err()
	r.Error(err)
	r.Equal(boom, err)
}

func Test_Context_Cancel(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	boom := fmt.Errorf("boom")

	outer, cancel := context.WithCancel(context.Background())
	defer cancel()

	inner, cancel := NewContext(outer)
	defer cancel()

	go func() {
		inner.SetErr(boom)
		inner.Cancel()
	}()

	<-inner.Done()

	err := inner.Err()
	r.Error(err)
	r.Equal(boom, err)
}
