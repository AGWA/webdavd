package main

import (
	"context"
	"errors"
	"os"

	"golang.org/x/net/webdav"
)

var errReadOnlyFilesystem = errors.New("read only filesystem")

type readOnlyFileSystem struct {
	webdav.FileSystem
}

func (fs readOnlyFileSystem) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	return errReadOnlyFilesystem
}

func (fs readOnlyFileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	if (flag&os.O_WRONLY) != 0 || (flag&os.O_RDWR) != 0 {
		return nil, errReadOnlyFilesystem
	}
	return fs.FileSystem.OpenFile(ctx, name, flag, perm)
}

func (fs readOnlyFileSystem) RemoveAll(ctx context.Context, name string) error {
	return errReadOnlyFilesystem
}

func (fs readOnlyFileSystem) Rename(ctx context.Context, oldName string, newName string) error {
	return errReadOnlyFilesystem
}
