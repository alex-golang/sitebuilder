// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	// FilePermission defines permissions for generated files.
	FilePermission = 0644

	// DirPermission defines permissions for generated directories.
	DirPermission = 0744
)

// ValidatePath ensures the given path is valid.
// This means it exists, and contains a few expected sub directories.
//
// It returns the absolute version of the path or an error.
func ValidatePath(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	err = dirExists(path)
	if err != nil {
		return "", err
	}

	err = dirExists(path, "posts")
	if err != nil {
		return "", err
	}

	err = dirExists(path, "static")
	if err != nil {
		return "", err
	}

	err = dirExists(path, "templates")
	if err != nil {
		return "", err
	}

	// Delete existing deploy directory.
	deploy := filepath.Join(path, "deploy")
	err = os.RemoveAll(deploy)
	if err != nil {
		return "", err
	}

	// Create new posts directory.
	err = os.MkdirAll(filepath.Join(deploy, "posts"), DirPermission)
	if err != nil {
		return "", err
	}

	// Create new tags directory.
	return path, os.MkdirAll(filepath.Join(deploy, "tags"), DirPermission)
}

// dirExists ensures the given path exists and that
// it is a directory.
func dirExists(paths ...string) error {
	path := strings.Join(paths, string(filepath.Separator))

	stat, err := os.Lstat(path)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return newError("Path is not a directory.")
	}

	return nil
}
