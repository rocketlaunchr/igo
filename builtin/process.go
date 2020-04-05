// Copyright 2019-20 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package builtin

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"

	"github.com/rocketlaunchr/igo/common"
	"github.com/rocketlaunchr/igo/file"
)

func Process(tempFile, sourceFile string) (rErr error) {
	defer func() {
		if err := recover(); err != nil {
			rErr = fmt.Errorf("%s: %v", sourceFile, err)
		}
	}()

	fset := token.NewFileSet()

	node, err := parser.ParseFile(fset, tempFile, nil, parser.AllErrors|parser.ParseComments|parser.DeclarationErrors)
	if err != nil {
		return err
	}

	mustFound := false

	astutil.Apply(node, pre(&mustFound), post)

	// If we found a "must" function then import "exported" package
	// TODO: required for future generics-compatible version of must.
	if mustFound {
		common.InsertImport(node, &[]string{"exported"}[0], "github.com/rocketlaunchr/igo/exported")
	}

	// Maintain the integrity of the comments
	for funcDecl, cg := range funcComments {
		for _, c := range cg.List {
			c.Slash = funcDecl.Pos() - 1
		}
	}
	funcComments = map[*ast.FuncDecl]*ast.CommentGroup{} // Reset for next file

	err = file.SaveFmtFile(tempFile, fset, node)
	if err != nil {
		return err
	}

	return nil

}
