// Copyright 2018 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package common

type UndoType int

const (
	DeferGo    UndoType = 0
	Fordefer            = 1
	FordeferGo          = 2
)

// Undo stores information about how to transform the converted code
// back to its original form.
// The Pos refers to the order of the identifier we are searching for
// in the transformed code.
type Undo struct {
	UndoType UndoType
	Pos      int // Optional
}
