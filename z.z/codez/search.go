package codez

import (
	"ezpkg.io/errorz"
)

type zKind int

const (
	zIdent zKind = iota + 1
	zExpr
	zStmt
	zDecl
	zFile
)

type zSearch = Search
type Search struct {
	pats []codePattern
	pkgs []string

	compiled bool
	errs     errorz.Errors
}

type codePattern struct {
	id  int
	pat string
}

func NewSearch(pattern string) *Search {
	return &Search{
		pats: []codePattern{{id: 0, pat: pattern}},
	}
}

func (s *Search) mustNotCompiled() {
	if s.compiled {
		panic("unexpected: search is already compiled, changes are not accepted")
	}
}

func (s *Search) Clone() *Search {
	return &Search{}
}

func (s *Search) Import(alias, pkg string) *Search {
	s.mustNotCompiled()
	return s
}

func (s *Search) WithIdent(name string) *SearchIdent {
	s.mustNotCompiled()
	return &SearchIdent{zSearch: s, zVar: name}
}

func (s *Search) WithExpr(name string) *SearchExpr {
	s.mustNotCompiled()
	return &SearchExpr{zSearch: s, zVar: name}
}

func (s *Search) WithStmt(name string) *SearchStmt {
	s.mustNotCompiled()
	return &SearchStmt{zSearch: s, zVar: name}
}

func (s *Search) AddPattern(id int, pattern string) *Search {
	s.mustNotCompiled()
	s.pats = append(s.pats, codePattern{id: id, pat: pattern})
	return s
}

func (s *Search) InPackages(pkgs ...string) *Search {
	s.mustNotCompiled()
	s.pkgs = append(s.pkgs, pkgs...)
	return s
}

func (s *Search) Exec(pkgs *Packages) (_ *SearchResult, errs errorz.Errors) {
	s.compiled = true
	return &SearchResult{}, nil
}

type SearchIdent struct {
	*zSearch

	zVar string
}
type SearchExpr struct {
	*zSearch

	zVar string
}
type SearchStmt struct {
	*zSearch

	zVar string
}

func (s *SearchIdent) IdentType(identTypes ...string) *SearchIdent {
	return s
}

func (s *SearchExpr) ExprType(exprTypes ...string) *SearchExpr {
	return s
}
