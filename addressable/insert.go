// Copyright 2018 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package addressable

import (
	"go/ast"
	"go/token"

	"github.com/rocketlaunchr/igo/common"
)

func insertSingleLine(typ string, val string) *ast.IndexExpr {
	return &ast.IndexExpr{
		Index: &ast.BasicLit{
			Kind:  token.INT,
			Value: "0",
		},
		X: &ast.CompositeLit{
			Type: &ast.ArrayType{
				Elt: &ast.Ident{
					Name: typ,
				},
			},
			Elts: []ast.Expr{
				&ast.Ident{
					Name: val,
				},
			},
		},
	}
}

func insertCallExpr(blockStmt *ast.BlockStmt, n int, val *ast.CallExpr) string {
	varName := common.RandSeq(randLength)
	identifier := ast.NewIdent(varName)
	assignment := &ast.AssignStmt{
		Lhs: []ast.Expr{identifier},
		Rhs: []ast.Expr{
			val,
		},
		Tok: token.DEFINE,
	}

	blockStmt.List = append(blockStmt.List[:n], append([]ast.Stmt{assignment}, blockStmt.List[n:]...)...)

	return varName
}

func insertBoolVar(blockStmt *ast.BlockStmt, n int, val string) string {
	varName := common.RandSeq(randLength)
	identifier := ast.NewIdent(varName)
	assignment := &ast.AssignStmt{
		Lhs: []ast.Expr{identifier},
		Rhs: []ast.Expr{
			&ast.Ident{
				Name: val,
			},
		},
		Tok: token.DEFINE,
	}

	blockStmt.List = append(blockStmt.List[:n], append([]ast.Stmt{assignment}, blockStmt.List[n:]...)...)

	return varName
}

func insertConstVar(blockStmt *ast.BlockStmt, kind token.Token, n int, val string) string {
	varName := common.RandSeq(randLength)
	identifier := ast.NewIdent(varName)
	assignment := &ast.AssignStmt{
		Lhs: []ast.Expr{identifier},
		Rhs: []ast.Expr{
			&ast.BasicLit{
				Kind:  kind,
				Value: val,
			},
		},
		Tok: token.DEFINE,
	}

	blockStmt.List = append(blockStmt.List[:n], append([]ast.Stmt{assignment}, blockStmt.List[n:]...)...)

	return varName
}

func replaceX(varName string) *ast.Ident {
	return &ast.Ident{
		Name: varName,
	}
}
