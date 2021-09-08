package cleo

import (
	"fmt"
	"sync"
)

type Router struct {
	*sync.RWMutex
	def    HandlerFn
	routes map[string]HandlerFn
}

func NewRouter() *Router {
	return &Router{
		RWMutex: &sync.RWMutex{},
		routes:  map[string]HandlerFn{},
	}
}

func (s *Router) SetDefault(fn HandlerFn) {
	s.Lock()
	s.def = fn
	s.Unlock()
}

func (s *Router) Set(name string, fn HandlerFn) {
	s.Lock()

	if s.routes == nil {
		s.routes = map[string]HandlerFn{}
	}

	s.routes[name] = fn
	s.Unlock()
}

func (s *Router) Default(rt *Runtime) error {
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

func (s *Router) Switch(rt *Runtime) error {
	if rt == nil {
		return fmt.Errorf("runtime is nil")
	}

	args := rt.Args
	if len(args) == 0 {
		return s.Default(rt)
	}

	s.Lock()
	if s.routes == nil {
		s.routes = map[string]HandlerFn{}
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
