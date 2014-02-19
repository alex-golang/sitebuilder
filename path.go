// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"os"
	"path/filepath"
	"strings"
)

// ValidatePath ensures the given path is valid.
// This means it exists, and contains a few expected sub directories.
func ValidatePath(path string) string {
	path, err := filepath.Abs(path)
	test(err, "Invalid path")

	validateDir(path)
	validateDir(path, "posts")
	validateDir(path, "static")
	validateDir(path, "templates")

	// Delete existing deploy directory.
	dir := filepath.Join(path, "deploy")
	err = os.RemoveAll(dir)
	test(err, "Create deploy directory")

	// Create fresh deploy directory.
	err = os.MkdirAll(dir, 0744)
	test(err, "Create deploy directory")
	return path
}

// validateDir ensures the given path exists and that
// it is a directory.
func validateDir(paths ...string) {
	path := strings.Join(paths, string(filepath.Separator))

	stat, err := os.Lstat(path)
	test(err, "Invalid path")

	// ...and that it is a directory.
	if !stat.IsDir() {
		fatal("Path is not a directory.")
	}
}
