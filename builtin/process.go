// Copyright 2019 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package builtin

import (
	"fmt"
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

	// If we found a "must" function then import exported package
	if mustFound {
		common.InsertImport(node, &[]string{"exported"}[0], "github.com/rocketlaunchr/igo/exported")
	}

	err = file.SaveFmtFile(tempFile, fset, node)
	if err != nil {
		return err
	}

	return nil

}
