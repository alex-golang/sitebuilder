// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	"os"
)

// fatal logs the given formatted error and exits the program.
func fatal(msg string, argv ...interface{}) {
	fmt.Fprintf(os.Stderr, msg, argv...)
	os.Exit(1)
}

// warn displays a formatted warning message.
func warn(msg string, argv ...interface{}) {
	fmt.Printf(msg, argv...)
}

// newError creates a new, formatted error.
func newError(msg string, argv ...interface{}) error {
	return fmt.Errorf(msg, argv...)
}

// check tests the given error value. If not nil, the error
// is displayed and the program exits.
func check(err error) {
	if err != nil {
		fatal("%v", err)
	}
}
