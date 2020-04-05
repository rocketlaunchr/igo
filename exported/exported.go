// Copyright 2019-20 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package exported

// MustOptional is an optional second argument to the "builtin" must function.
// It allows you to wrap an error before the function panics.
type MustOptional func(error) error
