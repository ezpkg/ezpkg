package codez

import (
	"go/ast"
)

type VisitContext struct {
}

func Walk(node ast.Node, fn func(cx *VisitContext, node ast.Node) bool) {
	v := zVisitor{
		fn: fn,
	}
	v.visit(node)
}

type zVisitFunc func(cx *VisitContext, node ast.Node) bool

type zVisitor struct {
	cx *VisitContext
	fn zVisitFunc
}

func (v *zVisitor) visit(node ast.Node) {
	switch x := node.(type) {
	case ast.Expr:
		v.visitExpr(x)
	case ast.Stmt:
		v.visitStmt(x)
	case ast.Decl:
		v.visitDecl(x)
	case ast.Spec:
		v.visitSpec(x)
	case ast.Node:
		v.visitOther(x)
	}
}
