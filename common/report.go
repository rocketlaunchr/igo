// Copyright 2018-20 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package common

import (
	"go/scanner"
	"os"
)

func Report(err error) (exitCode int) {
	scanner.PrintError(os.Stderr, err)
	exitCode = 2
	return
}
