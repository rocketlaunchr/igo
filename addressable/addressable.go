// Copyright 2018 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package addressable

import (
	"go/ast"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"

	"github.com/rocketlaunchr/igo/file"
)

const randLength = 15 // TODO: Make this configurable

var encountredBlocks blocks

func init() {
	encountredBlocks = blocks{}
}

func Process(tempFile string) error {

	fset := token.NewFileSet()

	// Don't parse comments because it interferes with inserting lines above code
	node, err := parser.ParseFile(fset, tempFile, nil, parser.AllErrors|parser.DeclarationErrors)
	if err != nil {
		return err
	}

	astutil.Apply(node, pre, post)

	err = file.SaveFmtFile(tempFile, fset, node)
	if err != nil {
		return err
	}

	return nil
}

func pre(c *astutil.Cursor) bool {

	currentNode := c.Node()

	switch n := currentNode.(type) {
	case *ast.BlockStmt:
		encountredBlocks.add(n)
	case *ast.AssignStmt, *ast.BadStmt, *ast.BranchStmt, *ast.DeclStmt, *ast.DeferStmt,
		*ast.EmptyStmt, *ast.ExprStmt, *ast.ForStmt, *ast.GoStmt, *ast.IfStmt, *ast.IncDecStmt,
		*ast.LabeledStmt, *ast.RangeStmt, *ast.ReturnStmt, *ast.SelectStmt, *ast.SendStmt,
		*ast.SwitchStmt, *ast.TypeSwitchStmt:
		if currentBlock := encountredBlocks.current(); currentBlock != nil {
			// Search in current block to see if this stmt is a direct child of parent "blockStmt"
			if idx, exists := currentBlock.lookup[n.(ast.Stmt)]; exists {
				encountredBlocks[len(encountredBlocks)-1].current = &[]int{idx}[0]
			}
		}
	case *ast.UnaryExpr:
		if n.Op == token.AND { // Address Operator

			switch t := n.X.(type) {
			case *ast.Ident:
				// Note: Assume "true" and "false" are not redefined from default boolean type
				if t.Name == "true" || t.Name == "false" {
					if currentBlock := encountredBlocks.current(); currentBlock != nil {
						n.X = insertSingleLine("bool", t.Name)
					}
				}
			case *ast.CallExpr:
				if currentBlock := encountredBlocks.current(); currentBlock != nil {
					varName := insertCallExpr(currentBlock.node, *currentBlock.current, t)
					n.X = replaceX(varName)
				}
			case *ast.BasicLit:
				switch t.Kind {
				case token.STRING:
					if currentBlock := encountredBlocks.current(); currentBlock != nil {
						n.X = insertSingleLine("string", t.Value)
					}
				case token.INT, token.FLOAT, token.IMAG, token.CHAR:
					if currentBlock := encountredBlocks.current(); currentBlock != nil {
						varName := insertConstVar(currentBlock.node, t.Kind, *currentBlock.current, t.Value)
						n.X = replaceX(varName)
					}
				}
			}
		}
	}

	return true
}

func post(c *astutil.Cursor) bool {

	currentNode := c.Node()

	switch currentNode.(type) {
	case *ast.BlockStmt:
		encountredBlocks.remove()
	}

	return true
}
