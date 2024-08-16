package codez

import (
	"go/ast"
	"regexp"
)

type Matcher interface {
	Match(node ast.Node) (bool, error)
}

type StringMatcher struct {
	Value  string
	Regexp *regexp.Regexp
}

type BoolMatcher struct {
	Value bool
}

type ExprMatcher struct {
}

type StmtMatcher struct {
}

type StmtListMatcher struct {
}

type DeclMatcher struct {
}

type ExprListMatcher struct {
}

type ChanDirMatcher struct {
	_ *ast.ChanDir
}

type ObjectMatcher struct{}

type FieldMatcher struct {
}

type FieldListMatcher struct {
}

type CommentGroupMatcher struct {
}

type SpecMatcher struct {
}

type SpecListMatcher struct {
}

func MatchIdent(name string) IdentMatcher {
	return IdentMatcher{
		Name: StringMatcher{Value: name},
	}
}
