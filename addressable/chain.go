// Copyright 2018-19 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package addressable

import (
	"go/ast"
)

type block struct {
	node    *ast.BlockStmt
	lookup  map[ast.Stmt]int // Maps each stmt in the block to its index
	current *int             // signifies which stmt in block we are walking through at the moment
}

type blocks []block

func (b *blocks) add(bs *ast.BlockStmt) {

	bl := block{
		node:   bs,
		lookup: map[ast.Stmt]int{},
	}

	for i, v := range bs.List {
		bl.lookup[v] = i
	}

	*b = append(*b, bl)
}

func (b *blocks) remove() {
	if len(*b) > 0 {
		*b = (*b)[:len(*b)-1]
	}
}

func (b *blocks) current() *block {
	if len(*b) == 0 {
		return nil
	}

	return &(*b)[len(*b)-1]
}
