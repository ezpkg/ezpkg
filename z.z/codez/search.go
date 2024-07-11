package codez

import (
	"ezpkg.io/errorz"
)

type zKind int

const (
	zIdent zKind = iota + 1
	zExpr
	zStmt
)

type zSearch = Search
type Search struct {
	pats []codePattern
	pkgs []string

	dirty bool
	errs  errorz.Errors
}

type codePattern struct {
	id  int
	pat string
}

func NewSearch(pattern string) *Search {
	return &Search{
		pats:  []codePattern{{id: 0, pat: pattern}},
		dirty: true,
	}
}

func (s *Search) Import(alias, pkg string) *Search {
	return s
}

func (s *Search) WithIdent(name string) *SearchIdent {
	return &SearchIdent{zSearch: s, zVar: name}
}

func (s *Search) WithExpr(name string) *SearchExpr {
	return &SearchExpr{zSearch: s, zVar: name}
}

func (s *Search) WithStmt(name string) *SearchStmt {
	return &SearchStmt{zSearch: s, zVar: name}
}

func (s *Search) AddPattern(id int, pattern string) *Search {
	s.dirty = true
	s.pats = append(s.pats, codePattern{id: id, pat: pattern})
	return s
}

func (s *Search) InPackages(pkgs ...string) *Search {
	s.dirty = true
	s.pkgs = append(s.pkgs, pkgs...)
	return s
}

func (s *Search) Validate(pkgs *Packages) (warns, errs errorz.Errors) {
	s.dirty = false
	return nil, nil
}

func (s *Search) Exec(pkgs *Packages) (*SearchResult, errorz.Errors) {
	if s.dirty {
		_, err := s.Validate(pkgs)
		if err != nil {
			return nil, err
		}
	}
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
