package cleo

import (
	"flag"
	"fmt"
	"sync"
)

type SwitchFn func(rt *Runtime) error

type Switcher struct {
	*flag.FlagSet
	*sync.RWMutex
	def    SwitchFn
	routes map[string]SwitchFn
}

func NewSwitcher() *Switcher {
	return &Switcher{
		FlagSet: flag.NewFlagSet("", flag.ExitOnError),
		RWMutex: &sync.RWMutex{},
		routes:  map[string]SwitchFn{},
	}
}

func (s *Switcher) SetDefault(fn SwitchFn) {
	s.Lock()
	s.def = fn
	s.Unlock()
}

func (s *Switcher) Set(name string, fn SwitchFn) {
	s.Lock()

	if s.routes == nil {
		s.routes = map[string]SwitchFn{}
	}

	s.routes[name] = fn
	s.Unlock()
}

func (s *Switcher) Default(rt *Runtime) error {
	if rt == nil {
		return fmt.Errorf("runtime is nil")
	}

	s.RLock()
	fn := s.def
	s.RUnlock()

	if fn == nil {
		return fmt.Errorf("no default runtime switch set")
	}

	return fn(rt)
}

func (s *Switcher) Switch(rt *Runtime) error {
	if rt == nil {
		return fmt.Errorf("runtime is nil")
	}

	args := rt.Args
	if len(args) == 0 {
		return s.Default(rt)
	}

	s.Lock()
	if s.routes == nil {
		s.routes = map[string]SwitchFn{}
	}
	s.Unlock()

	n := args[0]
	rt.Args = args[1:]

	fn, ok := s.routes[n]
	if !ok {
		return fmt.Errorf("unknown route %q [%s]", n, rt.Args)
	}

	return fn(rt)
}
