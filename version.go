// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	"runtime"
)

const (
	AppName         = "sitebuild"
	AppVersionMajor = 0
	AppVersionMinor = 1
)

func Version() string {
	return fmt.Sprintf("%s %d.%d (Go runtime %s).\nCopyright (c) 2010-2014, Jim Teeuwen.",
		AppName, AppVersionMajor, AppVersionMinor, runtime.Version())
}
