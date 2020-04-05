// Copyright 2019-20 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package common

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/rocketlaunchr/igo/file"
)

// RemoveComments will remove comments found inside functions to prevent interference
// when tree nodes are removed or manipulated.
func RemoveComments(tempFile string) error {

	fset := token.NewFileSet()

	// Don't parse comments because it interferes with inserting lines above code
	node, err := parser.ParseFile(fset, tempFile, nil, parser.AllErrors|parser.ParseComments|parser.DeclarationErrors)
	if err != nil {
		return err
	}

	type removeRange struct {
		Lbrace token.Pos // position of "{"
		Rbrace token.Pos // position of "}"
	}

	// remove all comments that fall in this range
	removeCommentsList := []removeRange{}

	ast.Inspect(node, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.FuncDecl:
			removeCommentsList = append(removeCommentsList, removeRange{Lbrace: n.Body.Lbrace, Rbrace: n.Body.Rbrace})
		case *ast.FuncLit:
			removeCommentsList = append(removeCommentsList, removeRange{Lbrace: n.Body.Lbrace, Rbrace: n.Body.Rbrace})
		}
		return true
	})

OUTER:
	for i := len(node.Comments) - 1; i >= 0; i-- {
		cg := node.Comments[i]
		comment := cg.List[0]

		// Check if comment is inside the removeCommentsList
		for _, v := range removeCommentsList {
			if comment.Slash >= v.Lbrace && comment.Slash <= v.Rbrace {
				// Remove entire comment group
				node.Comments = append(node.Comments[:i], node.Comments[i+1:]...)
				continue OUTER
			}
		}
	}

	err = file.SaveFmtFile(tempFile, fset, node)
	if err != nil {
		return err
	}

	return nil
}
