package codez

import (
	"go/ast"

	"ezpkg.io/errorz"
)

func Match(m Matcher, pkg *Package) (out []ast.Node, err error) {
	for _, f := range pkg.Syntax {
		ast.Inspect(f, func(node ast.Node) bool {
			ok, err0 := m.Match(node)
			switch {
			case err0 != nil:
				errorz.AppendTo(&err, err0)
				return false
			case ok:
				out = append(out, node)
				return false
			}
			return true
		})
	}
	return out, err
}
