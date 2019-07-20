// Copyright 2018 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package common

import (
	"fmt"
	"go/scanner"
	"go/token"
)

type state int

const (
	initial       state = 0
	deferFound          = 1
	fordeferFound       = 2
)

// FindIllegalStatements will search for these illegal statements:
// - "defer go"
// - "fordefer" and "fordefer go" (inside for loops)
// The function returns where the "go" statements are (before temporary removal)
// The function also returns where the "fordefer" statements are and whether they were
// next to a "go" statement.
func FindIllegalStatements(sourceFile string, sourceData []byte) ([]int, []Undo, []Undo, error) {

	goPos := []int{}

	// For undoing:
	undoDeferGoList := []Undo{}
	var defersFound int
	undoFordeferList := []Undo{}
	var fordefersFound int

	fs := token.NewFileSet()

	f := fs.AddFile(sourceFile, fs.Base(), len(sourceData))
	var s scanner.Scanner
	s.Init(f, sourceData, nil, 0)

	var state state
	var fordeferPos *int // When a fordefer is found, we record it's position

	for {
		pos, tok, l := s.Scan()
		if tok == token.EOF {
			break
		}

		if tok == token.DEFER {
			defersFound++

			state = deferFound
			fordeferPos = nil
		} else if tok == token.GO {
			switch state {
			case deferFound:
				undoDeferGoList = append(undoDeferGoList, Undo{DeferGo, defersFound})
				removeGoStmt(sourceData, int(pos))
				goPos = append(goPos, int(pos))
			case fordeferFound:
				// "fordefer _____ go" found
				undoFordeferList = append(undoFordeferList, Undo{FordeferGo, fordefersFound})
				replaceFordeferStmt(sourceData, *fordeferPos, int(pos)+2-*fordeferPos-len("fordefer"))
			}
			state = initial
			fordeferPos = nil
		} else if tok == token.IDENT && l == "fordefer" {
			// We found a fordefer statement
			fordefersFound++
			state = fordeferFound
			fordeferPos = &[]int{int(pos)}[0]
		} else {
			if state == fordeferFound {
				// Not a "fordefer go"
				undoFordeferList = append(undoFordeferList, Undo{Fordefer, fordefersFound})
				replaceFordeferStmt(sourceData, *fordeferPos, 0)
			}
			state = initial
			fordeferPos = nil
		}
	}

	if s.ErrorCount != 0 {
		return nil, nil, nil, fmt.Errorf("scanning failed. error count: %d", s.ErrorCount)
	}

	return goPos, undoDeferGoList, undoFordeferList, nil
}

// removeGoStmt replaces "go" statement with "  ".
func removeGoStmt(data []byte, pos int) {
	data[pos-1] = []byte(" ")[0]
	data[pos] = []byte(" ")[0]
}

// replaceFordeferStmt replaces "fordefer" or "fordefer go" with a 'fake'
// selector call ("f0.d5();").
func replaceFordeferStmt(data []byte, start int, gap int) {

	data[start-1] = []byte("f")[0]
	data[start] = []byte("0")[0]
	data[start+1] = []byte(".")[0]
	data[start+2] = []byte("d")[0]
	data[start+3] = []byte("5")[0]
	data[start+4] = []byte("(")[0]
	for i := 0; i < gap; i++ {
		data[start+5+i] = []byte(" ")[0]
	}
	data[start+5+gap] = []byte(")")[0]
	data[start+6+gap] = []byte(";")[0]
}
