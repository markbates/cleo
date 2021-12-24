package cleo

import (
	"io"
	"os"
	"sync"
)

type StdIO struct {
	in   io.Reader
	out  io.Writer
	err  io.Writer
	once sync.Once
}

func (s *StdIO) In() io.Reader {
	s.init()
	if s == nil {
		return nil
	}

	return s.in
}

func (s *StdIO) Out() io.Writer {
	s.init()
	if s == nil {
		return nil
	}

	return s.out
}

func (s *StdIO) Err() io.Writer {
	s.init()
	if s == nil {
		return nil
	}

	return s.err
}

func (s *StdIO) init() {
	if s == nil {
		return
	}

	s.once.Do(func() {
		if s.in == nil {
			s.in = os.Stdin
		}

		if s.out == nil {
			s.out = os.Stdout
		}

		if s.err == nil {
			s.err = os.Stderr
		}
	})

}

func WithIn(stdio *StdIO, r io.Reader) *StdIO {
	return &StdIO{
		in:  r,
		out: stdio.Out(),
		err: stdio.Err(),
	}
}

func WithOut(stdio *StdIO, w io.Writer) *StdIO {
	return &StdIO{
		in:  stdio.In(),
		out: w,
		err: stdio.Err(),
	}
}

func WithErr(stdio *StdIO, w io.Writer) *StdIO {
	return &StdIO{
		in:  stdio.In(),
		out: stdio.Out(),
		err: w,
	}
}
