// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

// CopyStatic copies static content from source to target
// directories, as-is.
func CopyStatic(root string) {
	src := filepath.Join(root, "static")
	dst := filepath.Join(root, "deploy")
	err := filepath.Walk(src, func(file string, stat os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		dst := strings.Replace(file, src, dst, 1)

		if stat.IsDir() {
			return os.MkdirAll(dst, DirPermission)
		}

		fs, err := os.Open(file)
		if err != nil {
			return err
		}

		defer fs.Close()

		fd, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, FilePermission)
		if err != nil {
			return err
		}

		defer fd.Close()

		_, err = io.Copy(fd, fs)
		return err
	})

	test(err, "Copy static content")
}
