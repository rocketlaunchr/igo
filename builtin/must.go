package builtin

import (
	"fmt"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"

	"github.com/rocketlaunchr/igo/common"
)

func replaceMustFunc(c *astutil.Cursor, arg1 *ast.CallExpr, arg2 ast.Expr) {

	// Assignment
	res1 := ast.NewIdent(common.RandSeq(6))
	res2 := ast.NewIdent(common.RandSeq(6))

	ass := &ast.AssignStmt{
		Lhs: []ast.Expr{
			res1,
			res2,
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			arg1,
		},
	}

	var panicCall ast.Expr
	if arg2 == nil {
		panicCall = res2
	} else {
		panicCall = &ast.CallExpr{
			Fun: arg2,
			Args: []ast.Expr{
				res2,
			},
		}
	}

	// Error to panic
	ifStmt := &ast.IfStmt{
		Cond: &ast.BinaryExpr{
			X:  res2,
			Op: token.NEQ,
			Y:  ast.NewIdent("nil"),
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: ast.NewIdent("panic"),
						Args: []ast.Expr{
							panicCall,
						},
					},
				},
			},
		},
	}

	insert := []ast.Stmt{ass, ifStmt}

	row := current().idx
	node := current().ref

	c.Replace(res1)

	switch n := node.(type) {
	case *ast.BlockStmt:
		n.List = append(n.List[:row], append(insert, n.List[row:]...)...)
	case *ast.CommClause:
		n.Body = append(n.Body[:row], append(insert, n.Body[row:]...)...)
	case *ast.CaseClause:
		n.Body = append(n.Body[:row], append(insert, n.Body[row:]...)...)
	}

	updateIdx(len(insert))
}

func isMustFunc(call *ast.CallExpr) (bool, *ast.CallExpr, ast.Expr) {

	// Check name
	t1, ok := call.Fun.(*ast.Ident)
	if !ok {
		return false, nil, nil
	}

	if t1.Name != "must" {
		return false, nil, nil
	}

	if len(call.Args) == 0 {
		panic("missing arguments to must")
	} else if len(call.Args) > 2 {
		panic("too many arguments to must")
	}

	// Is the first arg a function call?
	t2, ok := call.Args[0].(*ast.CallExpr)
	if !ok {
		panic(fmt.Sprintf("invalid operation: must (first argument has type %T, expecting function call)", call.Args[0]))
	}

	if len(call.Args) == 1 {
		return true, t2, nil
	} else {
		return true, t2, call.Args[1]
	}

}
