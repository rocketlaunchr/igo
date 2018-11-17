// Copyright 2018 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package fordefer

import (
	"reflect"
	"sync"
)

type tocall struct {
	fn        interface{}
	goroutine bool // Run in goroutine or not
}

// Stack is LIFO stack that stores closures.
type Stack struct {
	sync.Mutex
	stack []tocall
}

// NewStack creates a new Stack.
func NewStack() *Stack {
	return &Stack{
		stack: []tocall{},
	}
}

// Prepend inserts a closure to the beginning of the stack.
// Prepend is required for the stack to operate in LIFO mode.
func (s *Stack) Prepend(goroutine bool, fn interface{}) {
	s.Lock()
	defer s.Unlock()

	tc := tocall{
		fn:        fn,
		goroutine: goroutine,
	}

	s.stack = append([]tocall{tc}, s.stack...)
}

// Append inserts a closure to the end of the stack.
func (s *Stack) Append(goroutine bool, fn interface{}) {
	s.Lock()
	defer s.Unlock()

	tc := tocall{
		fn:        fn,
		goroutine: goroutine,
	}

	s.stack = append(s.stack, tc)
}

// Unwind executes all stored closures from the
// beginning to the end.
func (s *Stack) Unwind() {
	s.Lock()
	defer s.Unlock()

	for i := range s.stack {
		tc := s.stack[i]
		val := reflect.ValueOf(tc.fn)
		if tc.goroutine {
			go val.Call([]reflect.Value{})
		} else {
			val.Call([]reflect.Value{})
		}
	}

	// Reset stack back to original state
	s.stack = []tocall{}
}
