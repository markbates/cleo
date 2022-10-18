package cleo

import (
	"github.com/markbates/garlic"
)

func PWD(pwd string) (string, error) {
	return garlic.PWD(pwd)
}
