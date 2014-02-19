// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	"os"
)

// fatal writes the given formatted error to stderr and exits the program.
func fatal(msg string, argv ...interface{}) {
	msg = fmt.Sprintf(msg, argv...)
	fmt.Fprintf(os.Stderr, "%s\n", msg)
	os.Exit(1)
}

// test tests the given error value.
// If it is non-nil, it prints the given error and formatted
// message to stderr and then exits the program.
func test(err error, msg string, argv ...interface{}) {
	if err == nil {
		return
	}

	msg = fmt.Sprintf(msg, argv...)
	msg = fmt.Sprintf("%s: %%v\n", msg)

	fmt.Fprintf(os.Stderr, msg, err)
	os.Exit(1)
}
