package codez

import (
	"go/ast"
	"go/types"
)

type _MatchContext struct {
	pkg *Package
}

func newMatchContext(pkg *Package) *_MatchContext {
	return &_MatchContext{pkg: pkg}
}

func (cx *_MatchContext) getType(node ast.Node) types.Type {
	return nil
}
