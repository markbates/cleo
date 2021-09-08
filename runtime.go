package cleo

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/markbates/fsx"
)

type Runtimes []*Runtime

type Runtime struct {
	*flag.FlagSet
	Args   []string
	Cab    fs.FS
	Name   string
	Parent *Runtime
	Stderr io.Writer
	Stdin  io.Reader
	Stdout io.Writer
}

func (rt Runtime) Format(f fmt.State, verb rune) {
	fmt.Fprintf(f, "Runtime: %s", rt.Name)
	if len(rt.Args) > 0 {
		fmt.Fprintf(f, " [%s]", strings.Join(rt.Args, " "))
	}

	switch verb {
	case 'v':
		paths, _ := fsx.Paths(rt.Cab)
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
		FlagSet: rt.FlagSet,
		Name:    args[0],
		Parent:  rt,
		Stderr:  rt.Stderr,
		Stdin:   rt.Stdin,
		Stdout:  rt.Stdout,
	}

	return n, true
}

func (rt *Runtime) Parse(args []string) error {
	if err := rt.FlagSet.Parse(args); err != nil {
		return err
	}

	rt.Args = rt.FlagSet.Args()
	return nil
}

func RuntimeWithFlags(name string, args []string, flags *flag.FlagSet) (*Runtime, error) {
	if flags == nil {
		return nil, fmt.Errorf("flags can not be nil")
	}

	rt := &Runtime{
		Args:    args,
		FlagSet: flags,
		Name:    name,
		Stderr:  os.Stderr,
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
	}
	rt.SetOutput(rt.Stdout)

	if len(args) > 0 {
		rt.Name = args[0]
	}

	if err := rt.Parse(args); err != nil {
		return nil, err
	}

	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	cab, err := fsx.DirFS(pwd)
	if err != nil {
		return nil, err
	}

	rt.Cab = cab

	return rt, nil
}

func NewRuntime(name string, args []string) (*Runtime, error) {
	f := flag.NewFlagSet(name, flag.ExitOnError)
	return RuntimeWithFlags(name, args, f)
}
