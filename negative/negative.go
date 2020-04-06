// Copyright 2018-20 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package negative

import (
	"go/ast"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"

	"github.com/rocketlaunchr/igo/file"
)

func Process(tempFile string) error {

	fset := token.NewFileSet()

	node, err := parser.ParseFile(fset, tempFile, nil, parser.AllErrors|parser.ParseComments|parser.DeclarationErrors)
	if err != nil {
		return err
	}

	astutil.Apply(node, pre, nil)

	err = file.SaveFmtFile(tempFile, fset, node)
	if err != nil {
		return err
	}

	return nil
}

func pre(c *astutil.Cursor) bool {

	currentNode := c.Node()
	switch n := currentNode.(type) {
	case *ast.IndexExpr:

		// What is the symbol of the variable
		symbol, ok := n.X.(*ast.Ident)
		if !ok {
			break
		}

		// What is the index? Is it a negative int?
		if idx, ok := n.Index.(*ast.UnaryExpr); ok && idx.Op == token.SUB {
			if x, ok := idx.X.(*ast.BasicLit); ok && x.Kind == token.INT {
				// Replace index
				n.Index = replaceExpr(symbol, x.Value).X
			}
		}

	case *ast.SliceExpr:

		// What is the symbol of the variable
		symbol, ok := n.X.(*ast.Ident)
		if !ok {
			break
		}

		if n.Low != nil {
			if idx, ok := n.Low.(*ast.UnaryExpr); ok && idx.Op == token.SUB {
				if x, ok := idx.X.(*ast.BasicLit); ok && x.Kind == token.INT {
					// Replace index
					n.Low = replaceExpr(symbol, x.Value).X
				}
			}
		}

		if n.High != nil {
			if idx, ok := n.High.(*ast.UnaryExpr); ok && idx.Op == token.SUB {
				if x, ok := idx.X.(*ast.BasicLit); ok && x.Kind == token.INT {
					// Replace index
					n.High = replaceExpr(symbol, x.Value).X
				}
			}
		}

	}

	return true
}

func replaceExpr(symbol *ast.Ident, intVal string) *ast.ExprStmt {

	return &ast.ExprStmt{
		X: &ast.BinaryExpr{
			X: &ast.CallExpr{
				Fun: ast.NewIdent("len"),
				Args: []ast.Expr{
					symbol,
				},
			},
			Op: token.SUB,
			Y: &ast.BasicLit{
				Kind:  token.INT,
				Value: intVal,
			},
		},
	}

}
