package codez

import (
	"go/ast"
	"slices"
	"strings"

	"ezpkg.io/slicez"
)

type VisitContext struct {
	path    []string   // path to the current node
	nodes   []ast.Node // stack of nodes, not include the current node
	curNode ast.Node
}

func newVisitContext() *VisitContext {
	return &VisitContext{}
}

func (cx *VisitContext) Path() string {
	return strings.Join(cx.path, ".")
}
func (cx *VisitContext) Current() ast.Node {
	return cx.curNode
}
func (cx *VisitContext) Parent() ast.Node {
	return slicez.Last(cx.nodes)
}
func (cx *VisitContext) Ancestors() []ast.Node {
	return slices.Clone(cx.nodes)
}
func (cx *VisitContext) push(name string, node ast.Node) {
	cx.path = append(cx.path, name)
	cx.nodes = append(cx.nodes, cx.curNode)
	cx.curNode = node
}
func (cx *VisitContext) pop() {
	cx.curNode = slicez.Last(cx.nodes)
	cx.nodes = cx.nodes[:len(cx.nodes)-1]
	cx.path = cx.path[:len(cx.path)-1]
}

func Walk(node ast.Node, fn func(cx *VisitContext, node ast.Node) bool) {
	v := zVisitor{
		cx: newVisitContext(),
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
	v.cx.curNode = node // push the first node
	defer func() { v.cx.curNode = nil }()

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
