// Copyright (C) 2021 Andrew Ayer
//
// Permission is hereby granted, free of charge, to any person obtaining a
// copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
// THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
// OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.
//
// Except as contained in this notice, the name(s) of the above copyright
// holders shall not be used in advertising or otherwise to promote the
// sale, use or other dealings in this Software without prior written
// authorization.

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
