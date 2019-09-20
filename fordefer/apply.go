// Copyright 2018-19 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package fordefer

import (
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
)

func pre(forFound *bool) func(c *astutil.Cursor) bool {
	return func(c *astutil.Cursor) bool {
		currentNode := c.Node()
		parentNode := c.Parent()

		// Search for RangeStmt or ForStmt
		switch n := currentNode.(type) {
		case *ast.FuncLit, *ast.FuncDecl:
			// We encountered a "func() {}" or "func main() {"
			forLevels.add(n)
		case *ast.ExprStmt:
			// "fordefer" statement

			fordef := isExprStmtAForDefer(n)
			if fordef != other {

				forLoop := forloops.current()
				if forLoop != nil {
					forIdentifier := forLoop.id

					nextIdx := c.Index() + 1
					nextStmt := parentNode.(*ast.BlockStmt).List[nextIdx]

					if fordef == fordeferStd {
						morphFordefer(forIdentifier, n, nextStmt, parentNode, false, nextIdx)
					} else if fordef == fordeferGo {
						morphFordefer(forIdentifier, n, nextStmt, parentNode, true, nextIdx)
					}
				}
			}

		case *ast.SelectStmt, *ast.SwitchStmt:
			breakParents.add(n)
		case *ast.BranchStmt:
			switch n.Tok {
			case token.BREAK:
				if n.Label == nil && breakParents.current() == notForStmt {
					// If there is no label and it's inside a switch/select, then ignore it.
					return true
				}
				fallthrough
			case token.CONTINUE:
				if n.Label == nil {
					// no label for break or continue
					forLoop := forloops.current()
					if forLoop != nil {
						// Add an unwind above the break/continue statement
						c.InsertBefore(&ast.ExprStmt{X: stackUnwind(&ast.Ident{Name: forLoop.id})})
					}
				} else {
					// label scenario
					label := n.Label.Name

					forLoop := forloops.current()
					if forLoop != nil {
						cForInfo := lookup[forLoop.node]
						if cForInfo != nil {

							c.InsertBefore(&ast.ExprStmt{X: stackUnwind(&ast.Ident{Name: cForInfo.identifier})})
							cForInfo = cForInfo.parent

							for cForInfo != nil {
								// Inspect the parent
								if cForInfo.label == nil {
									c.InsertBefore(&ast.ExprStmt{X: stackUnwind(&ast.Ident{Name: cForInfo.identifier})})
									cForInfo = cForInfo.parent
								} else {
									c.InsertBefore(&ast.ExprStmt{X: stackUnwind(&ast.Ident{Name: cForInfo.identifier})})
									if *cForInfo.label == label {
										// We arrived at last "for" loop
										break
									} else {
										cForInfo = cForInfo.parent
									}
								}
							}
						}
					}
				}
			case token.GOTO:
				// TODO: Very complex. See: https://github.com/golang/go/issues/26058
			}
		case *ast.LabeledStmt:
			// Check if child is RangeStmt or ForStmt
			switch forStmt := n.Stmt.(type) {
			case *ast.RangeStmt, *ast.ForStmt:
				forID := insertStackAssignment(c) // Insert stack creation above "for" loop

				addForLoopToLookup(forStmt, forLevels.current(), forloops.current(), &n.Label.Name, forID.Name)

				// Add forloop to stack
				forloops.add(forID.Name, forStmt)

				breakParents.add(forStmt)

				// Insert Unwind @ end of "for" loop
				if rs, ok := forStmt.(*ast.RangeStmt); ok {
					rs.Body.List = append(rs.Body.List, &ast.ExprStmt{X: stackUnwind(forID)})
				}
				if fs, ok := forStmt.(*ast.ForStmt); ok {
					fs.Body.List = append(fs.Body.List, &ast.ExprStmt{X: stackUnwind(forID)})
				}

			}
		case *ast.RangeStmt, *ast.ForStmt:
			*forFound = true

			// Check if the for loop has a label above it
			labelAbove := checkIfLabelAbove(parentNode)
			if !labelAbove {
				forID := insertStackAssignment(c) // Insert stack creation above "for" loop

				addForLoopToLookup(n, forLevels.current(), forloops.current(), nil, forID.Name)

				// Add forloop to stack
				forloops.add(forID.Name, n)

				// Insert Unwind @ end of "for" loop
				if rs, ok := n.(*ast.RangeStmt); ok {
					rs.Body.List = append(rs.Body.List, &ast.ExprStmt{X: stackUnwind(forID)})
				}
				if fs, ok := n.(*ast.ForStmt); ok {
					fs.Body.List = append(fs.Body.List, &ast.ExprStmt{X: stackUnwind(forID)})
				}
			breakParents.add(n)

			}
		}

		return true
	}
}

func post(c *astutil.Cursor) bool {

	currentNode := c.Node()
	parentNode := c.Parent()

	switch n := currentNode.(type) {
	case *ast.FuncLit, *ast.FuncDecl:
		forLevels.remove()

	case *ast.SelectStmt, *ast.SwitchStmt:
		breakParents.remove()

	case *ast.LabeledStmt:
		// Check if child is RangeStmt or ForStmt
		switch n.Stmt.(type) {
		case *ast.RangeStmt, *ast.ForStmt:
			// Pop "for" loop from stack
			forloops.remove()

			breakParents.remove()
		}
	case *ast.RangeStmt, *ast.ForStmt:

		// Check if the for loop has a label above it
		labelAbove := checkIfLabelAbove(parentNode)
		if !labelAbove {
			// Pop "for" loop from stack
			forloops.remove()

			breakParents.remove()
		}
	}

	return true
}
