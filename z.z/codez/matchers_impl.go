package codez

import (
	"errors"
	"go/ast"
	"go/types"
	"regexp"
)

type NodeMatcher interface {
	Match(cx *MatchContext, node ast.Node) (bool, error)
}
type ListMatcher[NODE ast.Node] interface {
	Match(cx *MatchContext, nodes []NODE) (bool, error)
}
type ValueMatcher[V any] interface {
	Match(cx *MatchContext, value V) (bool, error)
}

type ExprMatcher interface {
	NodeMatcher
	MatchExpr(cx *MatchContext, expr ast.Expr) (bool, error)
}

type StmtMatcher interface {
	NodeMatcher
	MatchStmt(cx *MatchContext, stmt ast.Stmt) (bool, error)
}

type DeclMatcher interface {
	NodeMatcher
	MatchDecl(cx *MatchContext, decl ast.Decl) (bool, error)
}

type StringMatcher interface {
	Match(cx *MatchContext, value string) (bool, error)
}

type BoolMatcher interface {
	Match(cx *MatchContext, value bool) (bool, error)
}

type ExprListMatcher[NODE ast.Node] interface {
	Match(cx *MatchContext, nodes []NODE) (bool, error)
}

type StmtListMatcher[NODE ast.Node] interface {
	Match(cx *MatchContext, nodes []NODE) (bool, error)
}

type ChanDirMatcher interface {
	Match(cx *MatchContext, dir ast.ChanDir) (bool, error)
}

type ObjectMatcher struct{}

func (x ObjectMatcher) Match(cx *MatchContext, node ast.Node) (bool, error) { return true, nil }

type FieldMatcher interface {
}

type FieldListMatcher interface {
	Match(cx *MatchContext, node ast.Node) (bool, error)
}

type CommentGroupMatcher interface {
	Match(cx *MatchContext, node ast.Node) (bool, error)
}

type SpecMatcher interface {
	Match(cx *MatchContext, node ast.Node) (bool, error)
}

type SpecListMatcher[NODE ast.Spec] interface {
	Match(cx *MatchContext, nodes []NODE) (bool, error)
}

type zStringMatcher struct {
	Value  string
	Regexp *regexp.Regexp
}

func (m zStringMatcher) Match(cx *MatchContext, value string) (bool, error) {
	if m.Regexp != nil {
		return m.Regexp.MatchString(value), nil
	}
	return value == m.Value, nil
}

func MatchString(value string) StringMatcher {
	return zStringMatcher{Value: value}
}

func MatchRegexp(re *regexp.Regexp) StringMatcher {
	return zStringMatcher{Regexp: re}
}

func MatchRegexpStr(re string) StringMatcher {
	return zStringMatcher{Regexp: regexp.MustCompile(re)}
}

type zNilMatcher[M NodeMatcher] struct{ matchers []M }

func Nil[M NodeMatcher](matchers ...M) NodeMatcher {
	return zNilMatcher[M]{matchers: matchers}
}

func (m zNilMatcher[NodeMatcher]) Match(cx *MatchContext, node ast.Node) (bool, error) {
	if node == nil {
		return true, nil
	}
	if len(m.matchers) == 0 {
		return false, nil
	}
	return And(m.matchers...).Match(cx, node)
}

type zAndMatcher[M NodeMatcher] struct {
	Matchers []M
}

func And[M NodeMatcher](matchers ...M) NodeMatcher {
	return zAndMatcher[M]{Matchers: matchers}
}

func (m zAndMatcher[NodeMatcher]) Match(cx *MatchContext, node ast.Node) (bool, error) {
	if len(m.Matchers) == 0 {
		return false, errors.New("empty And matcher")
	}
	for _, matcher := range m.Matchers {
		ok, err := matcher.Match(cx, node)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

type zOrMatcher[M NodeMatcher] struct {
	Matchers []M
}

func Or[M NodeMatcher](matchers ...M) NodeMatcher {
	return zOrMatcher[M]{Matchers: matchers}
}

func (m zOrMatcher[NodeMatcher]) Match(cx *MatchContext, node ast.Node) (bool, error) {
	if len(m.Matchers) == 0 {
		return false, errors.New("empty Or matchers")
	}
	for _, matcher := range m.Matchers {
		ok, err := matcher.Match(cx, node)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

type zExprTypeMatcher struct {
	Type types.Type
}

func MatchExprType(t types.Type) NodeMatcher {
	return zExprTypeMatcher{Type: t}
}

func (m zExprTypeMatcher) Match(cx *MatchContext, node ast.Node) (bool, error) {
	expr, ok := node.(ast.Expr)
	if !ok {
		return false, nil
	}
	return cx.TypeOf(expr) == m.Type, nil
}

func MatchIdent(nameMatcher StringMatcher) IdentMatcher {
	return IdentMatcherB{
		Name: nameMatcher,
	}
}

func MatchIdentAny() IdentMatcher {
	return IdentMatcherB{
		Name: MatchRegexp(regexp.MustCompile(`.`)),
	}
}

func MatchBlockStmtAny() StmtMatcher {
	return MatchBlockStmtFunc(func(cx *MatchContext, stmt *ast.BlockStmt) (bool, error) {
		return true, nil
	})
}

func MatchBlockStmtFunc(fn func(cx *MatchContext, stmt *ast.BlockStmt) (bool, error)) StmtMatcher {
	return BlockStmtMatcherFunc(fn)
}

type BlockStmtMatcherFunc func(cx *MatchContext, stmt *ast.BlockStmt) (bool, error)

func (fn BlockStmtMatcherFunc) Match(cx *MatchContext, node ast.Node) (bool, error) {
	x, ok := node.(*ast.BlockStmt)
	if !ok {
		return false, nil
	}
	return fn(cx, x)
}
func (fn BlockStmtMatcherFunc) MatchStmt(cx *MatchContext, stmt ast.Stmt) (bool, error) {
	return fn(cx, stmt.(*ast.BlockStmt))
}

func MatchSelector(x ExprMatcher, sel IdentMatcher) SelectorExprMatcher {
	return SelectorExprMatcherB{X: x, Sel: sel}
}

func match(cx *MatchContext, ok bool, err error, m NodeMatcher, node ast.Node) (bool, error) {
	switch {
	case !ok:
		return false, err
	case m == nil && node == nil:
		return true, nil
	case m == nil || node == nil:
		return false, nil
	default:
		return m.Match(cx, node)
	}
}

func matchList[NODE ast.Node](cx *MatchContext, ok bool, err error, m ListMatcher[NODE], nodes []NODE) (bool, error) {
	if !ok {
		return false, err
	}
	return m.Match(cx, nodes)
}

func matchValue[V any](cx *MatchContext, ok bool, err error, m ValueMatcher[V], value V) (bool, error) {
	if !ok {
		return false, err
	}
	return m.Match(cx, value)
}
