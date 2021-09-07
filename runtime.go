package cli

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/markbates/fsx"
)

type Runtime struct {
	*flag.FlagSet
	Args   []string
	Cab    fs.FS
	Stderr io.Writer
	Stdin  io.Reader
	Stdout io.Writer
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
		Stderr:  os.Stderr,
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		FlagSet: flags,
	}
	rt.SetOutput(rt.Stdout)

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
