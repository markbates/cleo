package cleo

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
	"sync"
)

type Runtimes []*Runtime

type Runtime struct {
	*sync.RWMutex
	Args   []string
	Cab    fs.FS
	Name   string
	Parent *Runtime
	Stderr io.Writer
	Stdin  io.Reader
	Stdout io.Writer
	env    map[string]string
}

func (rt *Runtime) Setenv(key string, val string) {
	rt.Lock()
	defer rt.Unlock()
	if rt.env == nil {
		rt.env = map[string]string{}
	}
	rt.env[key] = val
}

func (rt *Runtime) Getenv(key string) (string, bool) {
	rt.RLock()
	defer rt.RUnlock()

	if rt.env == nil {
		return "", false
	}

	s, ok := rt.env[key]
	return s, ok
}

func (rt *Runtime) Format(f fmt.State, verb rune) {
	rt.RLock()
	defer rt.RUnlock()

	name := rt.Name
	if len(name) == 0 {
		name = "<empty>"
	}

	fmt.Fprintf(f, "Runtime: %s", name)
	if len(rt.Args) > 0 {
		fmt.Fprintf(f, " [%s]", strings.Join(rt.Args, " "))
	}

	switch verb {
	case 'v':
		paths, _ := Paths(rt.Cab)
		if len(paths) == 0 {
			return
		}

		fmt.Fprintf(f, "\n  Files:\n")
		for _, path := range paths {
			fmt.Fprintf(f, "    %s\n", path)
		}
	}
}

func (rt *Runtime) Next() (*Runtime, bool) {
	args := rt.Args
	if len(args) == 0 {
		return nil, false
	}

	n := &Runtime{
		Args:    args[1:],
		Cab:     rt.Cab,
		env:     rt.env,
		Name:    args[0],
		Parent:  rt,
		RWMutex: &sync.RWMutex{},
		Stderr:  rt.Stderr,
		Stdin:   rt.Stdin,
		Stdout:  rt.Stdout,
	}

	return n, true
}

// func (rt *Runtime) Parse(args []string) error {
// 	if err := rt.FlagSet.Parse(args); err != nil {
// 		return err
// 	}

// 	rt.Args = rt.FlagSet.Args()
// 	return nil
// }

func NewRuntime(name string, args []string) (*Runtime, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	env := map[string]string{}
	for _, pair := range os.Environ() {
		split := strings.Split(pair, "=")

		var key string
		var val string

		if len(split) > 0 {
			key = split[0]
		}

		if len(split) > 1 {
			val = split[1]
		}

		env[key] = val
	}

	rt := &Runtime{
		Args:    args,
		Cab:     os.DirFS(pwd),
		env:     env,
		Name:    name,
		RWMutex: &sync.RWMutex{},
		Stderr:  os.Stderr,
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
	}

	if len(rt.Name) == 0 && len(args) > 0 {
		rt.Name = args[0]
	}

	return rt, nil
}
