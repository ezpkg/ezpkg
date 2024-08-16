package codez

import (
	"go/ast"
)

type Matcher interface {
	Match(node ast.Node) (bool, error)
}

type ExprMatcher interface {
	Matcher
	MatchExpr(expr ast.Expr) (bool, error)
}

type StmtMatcher interface {
	Matcher
	MatchStmt(stmt ast.Stmt) (bool, error)
}

type DeclMatcher interface {
	Matcher
	MatchDecl(decl ast.Decl) (bool, error)
}

type StringMatcher interface {
}

type BoolMatcher interface {
}

type ExprListMatcher interface {
}

type StmtListMatcher interface {
}

type ChanDirMatcher interface {
}

type ObjectMatcher interface{}

type FieldMatcher interface {
}

type FieldListMatcher interface {
}

type CommentGroupMatcher interface {
}

type SpecMatcher interface {
}

type SpecListMatcher interface {
}

func MatchIdent(name string) IdentMatcher {
	panic("todo")
}

func MatchSelector(x ExprMatcher, sel IdentMatcher) SelectorExprMatcher {
	panic("todo")
}
