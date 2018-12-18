// Copyright 2018 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package addressable

import (
	"go/ast"
)

type blocks []*ast.BlockStmt

func (b *blocks) add(bs *ast.BlockStmt) {
	*b = append(*b, bs)
}

func (b *blocks) remove() {
	if len(*b) > 0 {
		*b = (*b)[:len(*b)-1]
	}
}

func (b *blocks) current() *ast.BlockStmt {
	if len(*b) == 0 {
		return nil
	}

	return &(*b)[len(*b)-1]
}
