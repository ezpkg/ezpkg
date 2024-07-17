package codez

import (
	"go/ast"
)

func Match(m Matcher, pkg *Package) (out []ast.Node) {
	for _, f := range pkg.Syntax {
		ast.Inspect(f, func(node ast.Node) bool {
			if m.Match(node) {
				out = append(out, node)
				return false
			}
			return true
		})
	}
	return out
}
