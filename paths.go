package cleo

import (
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
)

func Paths(cab fs.FS) ([]string, error) {
	var paths []string

	err := fs.WalkDir(cab, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		base := filepath.Base(path)
		if strings.HasPrefix(base, ".") {
			return nil
		}

		paths = append(paths, path)
		return nil
	})

	sort.Strings(paths)
	return paths, err
}
