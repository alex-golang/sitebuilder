// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"html/template"
	"os"
	"path/filepath"
)

// LoadTemplates loads a template set from the given directory.
func LoadTemplates(path string) *template.Template {
	path = filepath.Join(path, "templates")
	fd, err := os.Open(path)
	test(err, "Load templates")

	defer fd.Close()

	files, err := fd.Readdirnames(-1)
	test(err, "Load templates")

	for i := range files {
		files[i] = filepath.Join(path, files[i])
	}

	t, err := template.ParseFiles(files...)
	test(err, "Load templates")
	return t
}
