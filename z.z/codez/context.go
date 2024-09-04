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

func (cx *MatchContext) TypeOf(expr ast.Expr) types.Type {
	return cx.pkgs.TypeOf(expr)
}
