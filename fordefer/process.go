// Copyright 2018-20 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package fordefer

import (
	"go/ast"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"

	"github.com/rocketlaunchr/igo/common"
	"github.com/rocketlaunchr/igo/file"
)

const alias = "fordefer"
const randLength = 15 // TODO: Make this configurable

var (
	lookup       map[ast.Node]*forInfo // Key is "for" loop node
	forLevels    levels
	forloops     forLoopStack
	breakParents breakParentStack
)

func Process(tempFile string) error {

	lookup = map[ast.Node]*forInfo{}
	forLevels = levels{}
	forloops = forLoopStack{}
	breakParents = breakParentStack{}

	fset := token.NewFileSet()

	node, err := parser.ParseFile(fset, tempFile, nil, parser.AllErrors|parser.ParseComments|parser.DeclarationErrors)
	if err != nil {
		return err
	}

	astutil.Apply(node, pre, post)

	idsToRemove := map[string]struct{}{}
	forFound := false

	for _, fi := range lookup {
		if fi.count == 0 {
			// Remove stack creation
			idsToRemove[fi.identifier] = struct{}{}
		} else {
			forFound = true
		}
	}

	// If we found a "for" statement then import fordefer package
	if forFound {
		common.InsertImport(node, &[]string{alias}[0], "github.com/rocketlaunchr/igo/stack")
	}

	astutil.Apply(node, cleanse(idsToRemove), nil)

	err = file.SaveFmtFile(tempFile, fset, node)
	if err != nil {
		return err
	}

	return nil
}

// checkIfLabelAbove returns whether there is a label above "for" statement
func checkIfLabelAbove(parent ast.Node) bool {
	_, ok := parent.(*ast.LabeledStmt)
	if ok {
		return true
	}
	return false
}

func stackUnwind(identifier *ast.Ident) *ast.CallExpr {
	unwind := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X: identifier,
			Sel: &ast.Ident{
				Name: "Unwind",
			},
		},
	}
	return unwind
}

func insertStackAssignment(c *astutil.Cursor) *ast.Ident {

	identifier := ast.NewIdent(common.RandSeq(randLength))
	newStack := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X: &ast.Ident{
				Name: alias,
			},
			Sel: &ast.Ident{
				Name: "NewStack",
			},
		},
		Args: []ast.Expr{
			&ast.Ident{
				Name: "true",
			},
		},
	}
	assignment := &ast.AssignStmt{Lhs: []ast.Expr{identifier}, Rhs: []ast.Expr{newStack}, Tok: token.DEFINE}
	// <ident> := fordefer.NewStack()
	c.InsertBefore(assignment)
	// defer <ident>.Unwind()
	c.InsertBefore(&ast.DeferStmt{
		Call: stackUnwind(identifier),
	})

	return identifier
}
