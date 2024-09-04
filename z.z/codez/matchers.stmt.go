//go:build !genz

// Code generated by genz codez-matchers. DO NOT EDIT.

package codez

import (
	ast "go/ast"
	token "go/token"
)

// AssignStmt
type AssignStmtMatcherB struct {
	_ *ast.AssignStmt

	Lhs    ExprListMatcher[ast.Expr]
	TokPos token.Pos
	Tok    token.Token
	Rhs    ExprListMatcher[ast.Expr]
}

func (m AssignStmtMatcherB) MatchStmt(cx *_MatchContext, node ast.Stmt) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m AssignStmtMatcherB) Match(cx *_MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.AssignStmt)
	if !ok {
		return false, nil
	}
	ok, err = matchList(cx, ok, err, m.Lhs, x.Lhs)
	ok, err = matchList(cx, ok, err, m.Rhs, x.Rhs)
	return ok, err
}

// BlockStmt
type BlockStmtMatcherB struct {
	_ *ast.BlockStmt

	Lbrace token.Pos
	List   StmtListMatcher[ast.Stmt]
	Rbrace token.Pos
}

func (m BlockStmtMatcherB) MatchStmt(cx *_MatchContext, node ast.Stmt) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m BlockStmtMatcherB) Match(cx *_MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.BlockStmt)
	if !ok {
		return false, nil
	}
	ok, err = matchList(cx, ok, err, m.List, x.List)
	return ok, err
}

// BranchStmt
type BranchStmtMatcherB struct {
	_ *ast.BranchStmt

	TokPos token.Pos
	Tok    token.Token
	Label  IdentMatcher
}

func (m BranchStmtMatcherB) MatchStmt(cx *_MatchContext, node ast.Stmt) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m BranchStmtMatcherB) Match(cx *_MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.BranchStmt)
	if !ok {
		return false, nil
	}
	ok, err = match(cx, ok, err, m.Label, x.Label)
	return ok, err
}

// CaseClause
type CaseClauseMatcherB struct {
	_ *ast.CaseClause

	Case  token.Pos
	List  ExprListMatcher[ast.Expr]
	Colon token.Pos
	Body  StmtListMatcher[ast.Stmt]
}

func (m CaseClauseMatcherB) MatchStmt(cx *_MatchContext, node ast.Stmt) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m CaseClauseMatcherB) Match(cx *_MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.CaseClause)
	if !ok {
		return false, nil
	}
	ok, err = matchList(cx, ok, err, m.List, x.List)
	ok, err = matchList(cx, ok, err, m.Body, x.Body)
	return ok, err
}

// CommClause
type CommClauseMatcherB struct {
	_ *ast.CommClause

	Case  token.Pos
	Comm  StmtMatcher
	Colon token.Pos
	Body  StmtListMatcher[ast.Stmt]
}

func (m CommClauseMatcherB) MatchStmt(cx *_MatchContext, node ast.Stmt) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m CommClauseMatcherB) Match(cx *_MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.CommClause)
	if !ok {
		return false, nil
	}
	ok, err = match(cx, ok, err, m.Comm, x.Comm)
	ok, err = matchList(cx, ok, err, m.Body, x.Body)
	return ok, err
}

// DeclStmt
type DeclStmtMatcherB struct {
	_ *ast.DeclStmt

	Decl DeclMatcher
}

func (m DeclStmtMatcherB) MatchStmt(cx *_MatchContext, node ast.Stmt) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m DeclStmtMatcherB) Match(cx *_MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.DeclStmt)
	if !ok {
		return false, nil
	}
	ok, err = match(cx, ok, err, m.Decl, x.Decl)
	return ok, err
}

// DeferStmt
type DeferStmtMatcherB struct {
	_ *ast.DeferStmt

	Defer token.Pos
	Call  CallExprMatcher
}

func (m DeferStmtMatcherB) MatchStmt(cx *_MatchContext, node ast.Stmt) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m DeferStmtMatcherB) Match(cx *_MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.DeferStmt)
	if !ok {
		return false, nil
	}
	ok, err = match(cx, ok, err, m.Call, x.Call)
	return ok, err
}

// EmptyStmt
type EmptyStmtMatcherB struct {
	_ *ast.EmptyStmt

	Semicolon token.Pos
	Implicit  BoolMatcher
}

func (m EmptyStmtMatcherB) MatchStmt(cx *_MatchContext, node ast.Stmt) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m EmptyStmtMatcherB) Match(cx *_MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.EmptyStmt)
	if !ok {
		return false, nil
	}
	ok, err = matchValue(cx, ok, err, m.Implicit, x.Implicit)
	return ok, err
}

// ExprStmt
type ExprStmtMatcherB struct {
	_ *ast.ExprStmt

	X ExprMatcher
}

func (m ExprStmtMatcherB) MatchStmt(cx *_MatchContext, node ast.Stmt) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m ExprStmtMatcherB) Match(cx *_MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.ExprStmt)
	if !ok {
		return false, nil
	}
	ok, err = match(cx, ok, err, m.X, x.X)
	return ok, err
}

// ForStmt
type ForStmtMatcherB struct {
	_ *ast.ForStmt

	For  token.Pos
	Init StmtMatcher
	Cond ExprMatcher
	Post StmtMatcher
	Body BlockStmtMatcher
}

func (m ForStmtMatcherB) MatchStmt(cx *_MatchContext, node ast.Stmt) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m ForStmtMatcherB) Match(cx *_MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.ForStmt)
	if !ok {
		return false, nil
	}
	ok, err = match(cx, ok, err, m.Init, x.Init)
	ok, err = match(cx, ok, err, m.Cond, x.Cond)
	ok, err = match(cx, ok, err, m.Post, x.Post)
	ok, err = match(cx, ok, err, m.Body, x.Body)
	return ok, err
}

// GoStmt
type GoStmtMatcherB struct {
	_ *ast.GoStmt

	Go   token.Pos
	Call CallExprMatcher
}

func (m GoStmtMatcherB) MatchStmt(cx *_MatchContext, node ast.Stmt) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m GoStmtMatcherB) Match(cx *_MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.GoStmt)
	if !ok {
		return false, nil
	}
	ok, err = match(cx, ok, err, m.Call, x.Call)
	return ok, err
}

// IfStmt
type IfStmtMatcherB struct {
	_ *ast.IfStmt

	If   token.Pos
	Init StmtMatcher
	Cond ExprMatcher
	Body BlockStmtMatcher
	Else StmtMatcher
}

func (m IfStmtMatcherB) MatchStmt(cx *_MatchContext, node ast.Stmt) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m IfStmtMatcherB) Match(cx *_MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.IfStmt)
	if !ok {
		return false, nil
	}
	ok, err = match(cx, ok, err, m.Init, x.Init)
	ok, err = match(cx, ok, err, m.Cond, x.Cond)
	ok, err = match(cx, ok, err, m.Body, x.Body)
	ok, err = match(cx, ok, err, m.Else, x.Else)
	return ok, err
}

// IncDecStmt
type IncDecStmtMatcherB struct {
	_ *ast.IncDecStmt

	X      ExprMatcher
	TokPos token.Pos
	Tok    token.Token
}

func (m IncDecStmtMatcherB) MatchStmt(cx *_MatchContext, node ast.Stmt) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m IncDecStmtMatcherB) Match(cx *_MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.IncDecStmt)
	if !ok {
		return false, nil
	}
	ok, err = match(cx, ok, err, m.X, x.X)
	return ok, err
}

// LabeledStmt
type LabeledStmtMatcherB struct {
	_ *ast.LabeledStmt

	Label IdentMatcher
	Colon token.Pos
	Stmt  StmtMatcher
}

func (m LabeledStmtMatcherB) MatchStmt(cx *_MatchContext, node ast.Stmt) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m LabeledStmtMatcherB) Match(cx *_MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.LabeledStmt)
	if !ok {
		return false, nil
	}
	ok, err = match(cx, ok, err, m.Label, x.Label)
	ok, err = match(cx, ok, err, m.Stmt, x.Stmt)
	return ok, err
}

// RangeStmt
type RangeStmtMatcherB struct {
	_ *ast.RangeStmt

	For    token.Pos
	Key    ExprMatcher
	Value  ExprMatcher
	TokPos token.Pos
	Tok    token.Token
	Range  token.Pos
	X      ExprMatcher
	Body   BlockStmtMatcher
}

func (m RangeStmtMatcherB) MatchStmt(cx *_MatchContext, node ast.Stmt) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m RangeStmtMatcherB) Match(cx *_MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.RangeStmt)
	if !ok {
		return false, nil
	}
	ok, err = match(cx, ok, err, m.Key, x.Key)
	ok, err = match(cx, ok, err, m.Value, x.Value)
	ok, err = match(cx, ok, err, m.X, x.X)
	ok, err = match(cx, ok, err, m.Body, x.Body)
	return ok, err
}

// ReturnStmt
type ReturnStmtMatcherB struct {
	_ *ast.ReturnStmt

	Return  token.Pos
	Results ExprListMatcher[ast.Expr]
}

func (m ReturnStmtMatcherB) MatchStmt(cx *_MatchContext, node ast.Stmt) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m ReturnStmtMatcherB) Match(cx *_MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.ReturnStmt)
	if !ok {
		return false, nil
	}
	ok, err = matchList(cx, ok, err, m.Results, x.Results)
	return ok, err
}

// SelectStmt
type SelectStmtMatcherB struct {
	_ *ast.SelectStmt

	Select token.Pos
	Body   BlockStmtMatcher
}

func (m SelectStmtMatcherB) MatchStmt(cx *_MatchContext, node ast.Stmt) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m SelectStmtMatcherB) Match(cx *_MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.SelectStmt)
	if !ok {
		return false, nil
	}
	ok, err = match(cx, ok, err, m.Body, x.Body)
	return ok, err
}

// SendStmt
type SendStmtMatcherB struct {
	_ *ast.SendStmt

	Chan  ExprMatcher
	Arrow token.Pos
	Value ExprMatcher
}

func (m SendStmtMatcherB) MatchStmt(cx *_MatchContext, node ast.Stmt) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m SendStmtMatcherB) Match(cx *_MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.SendStmt)
	if !ok {
		return false, nil
	}
	ok, err = match(cx, ok, err, m.Chan, x.Chan)
	ok, err = match(cx, ok, err, m.Value, x.Value)
	return ok, err
}

// SwitchStmt
type SwitchStmtMatcherB struct {
	_ *ast.SwitchStmt

	Switch token.Pos
	Init   StmtMatcher
	Tag    ExprMatcher
	Body   BlockStmtMatcher
}

func (m SwitchStmtMatcherB) MatchStmt(cx *_MatchContext, node ast.Stmt) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m SwitchStmtMatcherB) Match(cx *_MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.SwitchStmt)
	if !ok {
		return false, nil
	}
	ok, err = match(cx, ok, err, m.Init, x.Init)
	ok, err = match(cx, ok, err, m.Tag, x.Tag)
	ok, err = match(cx, ok, err, m.Body, x.Body)
	return ok, err
}

// TypeSwitchStmt
type TypeSwitchStmtMatcherB struct {
	_ *ast.TypeSwitchStmt

	Switch token.Pos
	Init   StmtMatcher
	Assign StmtMatcher
	Body   BlockStmtMatcher
}

func (m TypeSwitchStmtMatcherB) MatchStmt(cx *_MatchContext, node ast.Stmt) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m TypeSwitchStmtMatcherB) Match(cx *_MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.TypeSwitchStmt)
	if !ok {
		return false, nil
	}
	ok, err = match(cx, ok, err, m.Init, x.Init)
	ok, err = match(cx, ok, err, m.Assign, x.Assign)
	ok, err = match(cx, ok, err, m.Body, x.Body)
	return ok, err
}
