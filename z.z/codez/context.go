package codez

import (
	"go/ast"
	"go/types"
)

type MatchContext struct {
	pkgs *Packages
}

func newMatchContext(pkgs *Packages) *MatchContext {
	return &MatchContext{pkgs: pkgs}
}

func (cx *MatchContext) GetType(node ast.Node) types.Type {
	return nil
}
