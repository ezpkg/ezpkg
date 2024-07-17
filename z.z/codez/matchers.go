package codez

import (
	"go/ast"
)

type Matcher interface {
	Match(n ast.Node) bool
}

type StmtMatcher struct {
}

type ExprMatcherI interface {
	Matcher
	_expr()
}

type _exprMatcher struct {
}

type StmtListMatcher struct {
	stmts []*StmtMatcher
}

type ExprListMatcher struct {
}

type StmtDeclMatcher struct {
	_ *ast.DeclStmt
}

type FuncDeclMatcher struct {
	_ *ast.FuncDecl

	Recv *FieldListMatcher
	Name *IdentMatcher
}

type FuncTypeMatcher struct {
	_ *ast.FuncType

	TypeParams *FieldListMatcher
	Params     *FieldListMatcher
	Results    *FieldListMatcher
}

type FieldListMatcher struct {
	_ *ast.FieldList

	Fields []*FieldMatcher
}

type FieldMatcher struct {
	_ *ast.Field

	Names []*IdentMatcher
	Type  ExprMatcherI
}

type SelectorExprMatcher struct {
	_ *ast.SelectorExpr
	_exprMatcher

	X   *IdentMatcher
	Sel ExprMatcherI
}

type IdentMatcher struct {
	_ *ast.Ident
	_exprMatcher

	Name string
}

func (m _exprMatcher) _expr() {}

func (m *IdentMatcher) Match(n ast.Node) bool {
	id, ok := n.(*ast.Ident)
	if !ok {
		return false
	}
	return id.Name == m.Name
}

func (m *SelectorExprMatcher) Match(n ast.Node) bool {
	sel, ok := n.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	return m.X.Match(sel.X) && m.Sel.Match(sel.Sel)
}
