package cleo

import "flag"

type Flaggable interface {
	Flags() *flag.FlagSet
}
