// Copyright 2018 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package defergo

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/rocketlaunchr/igo/file"
)

func Process(tempFile string, goPos []int) error {

	if len(goPos) != 0 {

		fset := token.NewFileSet()

		node, err := parser.ParseFile(fset, tempFile, nil, parser.AllErrors|parser.ParseComments|parser.DeclarationErrors)
		if err != nil {
			return err
		}

		deferStmts := []*ast.DeferStmt{}

		ast.Inspect(node, func(n ast.Node) bool {

			// Check if this is actually a "defer go" stmt
			deferStmt, ok := n.(*ast.DeferStmt)
			if ok {

				// Find position where "go" stmt could have been
				rangeStart := int(deferStmt.Pos()) + len("defer")
				rangeEnd := int(deferStmt.Call.Pos()) - 1

				for _, v := range goPos {
					if v >= rangeStart && v <= rangeEnd {
						deferStmts = append(deferStmts, deferStmt)
						continue
					}
				}
			}

			return true
		})

		// Modify the "defer go" statements
		for i := range deferStmts {
			deferStmt := deferStmts[i]

			goStmt := &ast.GoStmt{
				Call: deferStmt.Call,
			}

			deferStmt.Call = &ast.CallExpr{
				Fun: &ast.FuncLit{
					Type: &ast.FuncType{
						Params: &ast.FieldList{
							List: nil,
						},
						Results: nil,
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{goStmt},
					},
				},
				Args:     nil,
				Ellipsis: token.NoPos,
			}

		}

		err = file.SaveFmtFile(tempFile, fset, node)
		if err != nil {
			return err
		}
	}

	return nil
}
