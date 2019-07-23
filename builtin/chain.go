package builtin

import (
	"go/ast"
)

type tracker struct {
	ref ast.Node // Could be from BlockStmt, CaseClause or CommClause
	idx int
}

var chain []tracker = make([]tracker, 0)

func removeLast() {
	if len(chain) > 0 {
		chain = chain[:len(chain)-1]
	}
}

func current() tracker {
	return chain[len(chain)-1]
}

func updateIdx(row, delta int) {
	if len(chain) > 0 {
		chain[len(chain)-1].idx = chain[len(chain)-1].idx + delta
	}
}

// pop removes all items from the end that refer to the same slice of statements.
func pop(n ast.Node) {
	if len(chain) == 0 {
		return
	}

	for i := len(chain) - 1; i >= 0; i-- {
		tracker := chain[i]

		if tracker.ref == n {
			chain = chain[:len(chain)-1]
		} else {
			break
		}
	}
}
