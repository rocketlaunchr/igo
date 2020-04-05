// Copyright 2018-20 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package fordefer

import (
	"go/ast"
)

type fordefer int

const (
	other       fordefer = 0
	fordeferStd fordefer = 1
	fordeferGo  fordefer = 2
)

// isExprStmtAForDefer will return if we found a "fordefer" or "fordefer go"
// statement.
func isExprStmtAForDefer(expr *ast.ExprStmt) fordefer {

	// Check if X is a "*ast.CallExpr"
	call, ok := expr.X.(*ast.CallExpr)
	if !ok {
		return other
	}

	// Check if Fun is a SelectorExpr
	selEx, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return other
	}

	// Check if Sel.Name is "d5"
	if selEx.Sel.Name != "d5" {
		return other
	}

	// Check if X is a "*ast.Ident"
	idt, ok := selEx.X.(*ast.Ident)
	if !ok {
		return other
	}

	//Check if function receiver's Name is "f0"
	if idt.Name != "f0" {
		return other
	}

	// Check if "fordefer go" statement
	if call.Rparen == call.Lparen+1 {
		return fordeferStd
	}

	return fordeferGo
}

// morphFordefer will convert a "f0.d5();" marker into an
// actual fordefer statement.
func morphFordefer(forID string, fordeferStmt *ast.ExprStmt, nextStmt ast.Stmt, parentNode ast.Node, goroutine bool, nextStmtIdx int) {

	var goroutineStr string
	if goroutine {
		goroutineStr = "true"
	} else {
		goroutineStr = "false"
	}

	// Modify current node
	fordeferStmt.X.(*ast.CallExpr).Fun.(*ast.SelectorExpr).X.(*ast.Ident).Name = forID // identifier for for loop's stack
	fordeferStmt.X.(*ast.CallExpr).Fun.(*ast.SelectorExpr).Sel.Name = "Add"
	fordeferStmt.X.(*ast.CallExpr).Args = []ast.Expr{
		&ast.Ident{
			Name: goroutineStr,
		},
		&ast.FuncLit{
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: nil,
				},
				Results: nil,
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{nextStmt},
			},
		},
	}

	// Remove next node
	parentNode.(*ast.BlockStmt).List = append(parentNode.(*ast.BlockStmt).List[:nextStmtIdx], parentNode.(*ast.BlockStmt).List[nextStmtIdx+1:]...)
}

func addForLoopToLookup(forNode ast.Node, currentLevel *level, currentFor *forLoop, label *string, randIdentifier string) {

	// currentLevel := forLevels.current()
	// currentFor := forloops.current() // This is the prior "for" loop
	if currentFor != nil {
		cForInfo := lookup[currentFor.node]

		// Check if current level is same as prior "for" loop's level.
		if currentLevel.funcNode == cForInfo.level.funcNode {
			// prior "for" loop is parent
			lookup[forNode] = &forInfo{
				level:      currentLevel,
				identifier: randIdentifier,
				label:      label,
				parent:     cForInfo,
			}
		} else {
			lookup[forNode] = &forInfo{
				level:      currentLevel,
				identifier: randIdentifier,
				label:      label,
				parent:     nil,
			}
		}
	} else {
		// This is the first "for" loop encountered
		lookup[forNode] = &forInfo{
			level:      currentLevel,
			identifier: randIdentifier,
			label:      label,
			parent:     nil,
		}
	}
}
