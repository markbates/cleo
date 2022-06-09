package cleo

import "io/fs"

type FSable interface {
	FileSystem() fs.FS
}

type SetFSable interface {
	SetFileSystem(cab fs.FS)
}
