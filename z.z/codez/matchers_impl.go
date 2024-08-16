package codez

import (
	"go/ast"
	"regexp"
)

type NodeMatcher interface {
	Match(node ast.Node) (bool, error)
}
type ListMatcher[NODE ast.Node] interface {
	Match(nodes []NODE) (bool, error)
}
type ValueMatcher[V any] interface {
	Match(value V) (bool, error)
}

type ExprMatcher interface {
	NodeMatcher
	MatchExpr(expr ast.Expr) (bool, error)
}

type StmtMatcher interface {
	NodeMatcher
	MatchStmt(stmt ast.Stmt) (bool, error)
}

type DeclMatcher interface {
	NodeMatcher
	MatchDecl(decl ast.Decl) (bool, error)
}

type StringMatcher interface {
	Match(value string) (bool, error)
}

type BoolMatcher interface {
	Match(value bool) (bool, error)
}

type ExprListMatcher[NODE ast.Node] interface {
	Match(nodes []NODE) (bool, error)
}

type StmtListMatcher[NODE ast.Node] interface {
	Match(nodes []NODE) (bool, error)
}

type ChanDirMatcher interface {
	Match(dir ast.ChanDir) (bool, error)
}

type ObjectMatcher struct{}

func (x ObjectMatcher) Match(node ast.Node) (bool, error) { return true, nil }

type FieldMatcher interface {
}

type FieldListMatcher interface {
	Match(node ast.Node) (bool, error)
}

type CommentGroupMatcher interface {
	Match(node ast.Node) (bool, error)
}

type SpecMatcher interface {
	Match(node ast.Node) (bool, error)
}

type SpecListMatcher[NODE ast.Node] interface {
	Match(nodes []NODE) (bool, error)
}

type zStringMatcher struct {
	Value  string
	Regexp *regexp.Regexp
}

func (m zStringMatcher) Match(value string) (bool, error) {
	if m.Regexp != nil {
		return m.Regexp.MatchString(value), nil
	}
	return value == m.Value, nil
}

func MatchIdent(name string) IdentMatcher {
	return zIdentMatcher{
		Name: zStringMatcher{Value: name},
	}
}

func MatchSelector(x ExprMatcher, sel IdentMatcher) SelectorExprMatcher {
	return zSelectorExprMatcher{X: x, Sel: sel}
}

func match(ok bool, err error, m NodeMatcher, node ast.Node) (bool, error) {
	if !ok {
		return false, err
	}
	return m.Match(node)
}

func matchList[NODE ast.Node](ok bool, err error, m ListMatcher[NODE], nodes []NODE) (bool, error) {
	if !ok {
		return false, err
	}
	return m.Match(nodes)
}

func matchValue[V any](ok bool, err error, m ValueMatcher[V], value V) (bool, error) {
	if !ok {
		return false, err
	}
	return m.Match(value)
}
