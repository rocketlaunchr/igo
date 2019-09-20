package fordefer

import (
	"go/ast"
	"golang.org/x/tools/go/ast/astutil"
)

func cleanse(ids map[string]struct{}) func(c *astutil.Cursor) bool {
	return func(c *astutil.Cursor) bool {
		currentNode := c.Node()

		switch n := currentNode.(type) {
		case *ast.AssignStmt:
			if isFordeferNewStack(n, ids) {
				c.Delete()
				return true
			}
		case *ast.DeferStmt:
			callExpr := n.Call
			if isUnwind(callExpr, ids) {
				c.Delete()
				return true
			}
		case *ast.ExprStmt:
			if callExpr, ok := n.X.(*ast.CallExpr); ok {
				if isUnwind(callExpr, ids) {
					c.Delete()
					return true
				}
			}
		}

		return true
	}
}

// Searching for: <ident> := fordefer.NewStack()
// Return true to delete node.
func isFordeferNewStack(n *ast.AssignStmt, ids map[string]struct{}) bool {
	if len(n.Lhs) == 0 {
		return false
	}

	lhs := n.Lhs[0]

	if x, ok := lhs.(*ast.Ident); ok {
		name := x.Name
		if _, exists := ids[name]; exists {
			return true
		}
	}

	return false
}

// Searching for: defer <ident>.Unwind() or <ident>.Unwind()
// Return true to delete node.
func isUnwind(n *ast.CallExpr, ids map[string]struct{}) bool {
	fn := n.Fun

	if x, ok := fn.(*ast.SelectorExpr); ok {
		x := x.X
		if x, ok := x.(*ast.Ident); ok {
			name := x.Name
			if _, exists := ids[name]; exists {
				return true
			}
		}
	}

	return false
}
