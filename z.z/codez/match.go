package codez

import (
	"fmt"
	"go/ast"

	"ezpkg.io/errorz"
)

func Match(pkgs *Packages, m NodeMatcher, pkgPatterns ...string) (out []*NodeX, err error) {
	cx := newMatchContext(pkgs)
	for _, pkg := range pkgs.AllPackages(pkgPatterns...) {
		fmt.Println("-- pkg", pkg.PkgPath)
		for _, f := range pkg.Syntax {
			fmt.Println("-->", pkg.PkgPath, f.Name)
			Walk(f, func(vx *VisitContext, node ast.Node) bool {
				ok, err0 := m.Match(cx, node)
				switch {
				case err0 != nil:
					errorz.AppendTo(&err, err0)
					return false
				case ok:
					// TODO: add parent and replacer
					nodeX := newNodeX(node, nil, nil)
					out = append(out, nodeX)
					return false
				}
				return true
			})
		}
	}
	return out, err
}
