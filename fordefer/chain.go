// Copyright 2018-19 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package fordefer

import (
	"go/ast"
)

type forInfo struct {
	level      *level
	identifier string   // random identifier
	label      *string  // label above "for" loop
	parent     *forInfo // parent refers to
}

type forLoop struct {
	node ast.Node // *ast.RangeStmt or *ast.ForStmt
	id   string   // random identifier provided to "for" loop
}

// forLoopStack is used to keep track of ALL "for" loops we encounter.
type forLoopStack []forLoop

func (l *forLoopStack) add(id string, forNode ast.Node) {

	switch node := forNode.(type) {
	case *ast.RangeStmt, *ast.ForStmt:
		*l = append(*l, forLoop{
			node: node,
			id:   id,
		})
	default:
		panic("not a for stmt")
	}
}

func (l *forLoopStack) remove() {
	if len(*l) > 0 {
		*l = (*l)[:len(*l)-1]
	}
}

func (l *forLoopStack) current() *forLoop {
	if len(*l) == 0 {
		return nil
	}
	return &(*l)[len(*l)-1]
}

type level struct {
	funcNode ast.Node // *ast.FuncLit or *ast.FuncDecl
}

// levels is used to keep track of ALL function literals and function
// declarations. All "for" loops are "owned" by a function.
// When we encounter a new function all "for" loops inside that new
// function are owned by the new function.
type levels []level

func (l *levels) add(funcNode ast.Node) {

	switch node := funcNode.(type) {
	case *ast.FuncLit, *ast.FuncDecl:
		*l = append(*l, level{
			funcNode: node,
		})
	default:
		panic("not a func stmt")
	}

}

func (l *levels) remove() {
	if len(*l) > 0 {
		*l = (*l)[:len(*l)-1]
	}
}

func (l *levels) current() *level {
	if len(*l) == 0 {
		return nil
	}

	return &(*l)[len(*l)-1]
}

// breakParentStack is used to determine if an encountered "break" statement
// is immediately enclosed by a "for", "switch" or "select" statement.
type breakParentStack []ast.Node

func (l *breakParentStack) add(breakParent ast.Node) {

	switch breakParent.(type) {
	case *ast.RangeStmt, *ast.ForStmt, *ast.SwitchStmt, *ast.SelectStmt:
		*l = append(*l, breakParent)
	default:
		panic("not a for, switch or select statement")
	}
}

func (l *breakParentStack) remove() {
	if len(*l) > 0 {
		*l = (*l)[:len(*l)-1]
	}
}

type breakParent int

const (
	forStmt    breakParent = 0
	notForStmt breakParent = 1
)

func (l *breakParentStack) current() breakParent {
	if len(*l) == 0 {
		return notForStmt
	}
	x := (*l)[len(*l)-1]
	switch x.(type) {
	case *ast.RangeStmt, *ast.ForStmt:
		return forStmt
	default:
		return notForStmt
	}
}
