// Copyright 2019 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package common

import (
	"fmt"
	"go/ast"
	"go/token"
)

// InsertImport will add an import declaration to a file (if it doesn't exist already).
func InsertImport(node *ast.File, alias *string, path string) {

	path = fmt.Sprintf("\"%s\"", path)

	// Check if import already exists
	for _, impt := range node.Imports {

		// Check import path
		if impt.Path != nil {
			if impt.Path.Kind == token.STRING && impt.Path.Value == path {
				// Check alias
				if impt.Name == nil {
					if alias == nil {
						return
					}
				} else {
					if alias != nil {
						if impt.Name.Name == *alias {
							return
						}
					}
				}
			}
		}
	}

	// Import doesn't exist so add it
	importSpec := &ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: path}}
	if alias != nil {
		importSpec.Name = &ast.Ident{Name: *alias}
	}

	node.Imports = append([]*ast.ImportSpec{importSpec}, node.Imports...)

	node.Decls = append([]ast.Decl{&ast.GenDecl{
		TokPos: node.Package,
		Tok:    token.IMPORT,
		Specs:  []ast.Spec{importSpec},
	}}, node.Decls...)
}
