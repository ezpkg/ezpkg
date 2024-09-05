//go:build !genz

// Code generated by genz codez-matchers. DO NOT EDIT.

package codez

import (
	ast "go/ast"
	token "go/token"
)

// ArrayType
type ArrayTypeMatcherB struct {
	_ *ast.ArrayType

	Lbrack token.Pos
	Len    ExprMatcher
	Elt    ExprMatcher
}

func (m ArrayTypeMatcherB) MatchExpr(cx *MatchContext, node ast.Expr) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m ArrayTypeMatcherB) Match(cx *MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.ArrayType)
	if !ok {
		return false, nil
	}
	ok, err = match(cx, ok, err, m.Len, x.Len)
	ok, err = match(cx, ok, err, m.Elt, x.Elt)
	return ok, err
}

// BasicLit
type BasicLitMatcherB struct {
	_ *ast.BasicLit

	ValuePos token.Pos
	Kind     token.Token
	Value    StringMatcher
}

func (m BasicLitMatcherB) MatchExpr(cx *MatchContext, node ast.Expr) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m BasicLitMatcherB) Match(cx *MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.BasicLit)
	if !ok {
		return false, nil
	}
	ok, err = matchValue(cx, ok, err, m.Value, x.Value)
	return ok, err
}

// BinaryExpr
type BinaryExprMatcherB struct {
	_ *ast.BinaryExpr

	X     ExprMatcher
	OpPos token.Pos
	Op    token.Token
	Y     ExprMatcher
}

func (m BinaryExprMatcherB) MatchExpr(cx *MatchContext, node ast.Expr) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m BinaryExprMatcherB) Match(cx *MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.BinaryExpr)
	if !ok {
		return false, nil
	}
	ok, err = match(cx, ok, err, m.X, x.X)
	ok, err = match(cx, ok, err, m.Y, x.Y)
	return ok, err
}

// CallExpr
type CallExprMatcherB struct {
	_ *ast.CallExpr

	Fun      ExprMatcher
	Lparen   token.Pos
	Args     ListMatcher[ast.Expr]
	Ellipsis token.Pos
	Rparen   token.Pos
}

func (m CallExprMatcherB) MatchExpr(cx *MatchContext, node ast.Expr) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m CallExprMatcherB) Match(cx *MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.CallExpr)
	if !ok {
		return false, nil
	}
	ok, err = match(cx, ok, err, m.Fun, x.Fun)
	ok, err = matchList(cx, ok, err, m.Args, x.Args)
	return ok, err
}

// ChanType
type ChanTypeMatcherB struct {
	_ *ast.ChanType

	Begin token.Pos
	Arrow token.Pos
	Dir   ChanDirMatcher
	Value ExprMatcher
}

func (m ChanTypeMatcherB) MatchExpr(cx *MatchContext, node ast.Expr) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m ChanTypeMatcherB) Match(cx *MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.ChanType)
	if !ok {
		return false, nil
	}
	ok, err = matchValue(cx, ok, err, m.Dir, x.Dir)
	ok, err = match(cx, ok, err, m.Value, x.Value)
	return ok, err
}

// CompositeLit
type CompositeLitMatcherB struct {
	_ *ast.CompositeLit

	Type       ExprMatcher
	Lbrace     token.Pos
	Elts       ListMatcher[ast.Expr]
	Rbrace     token.Pos
	Incomplete BoolMatcher
}

func (m CompositeLitMatcherB) MatchExpr(cx *MatchContext, node ast.Expr) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m CompositeLitMatcherB) Match(cx *MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.CompositeLit)
	if !ok {
		return false, nil
	}
	ok, err = match(cx, ok, err, m.Type, x.Type)
	ok, err = matchList(cx, ok, err, m.Elts, x.Elts)
	ok, err = matchValue(cx, ok, err, m.Incomplete, x.Incomplete)
	return ok, err
}

// Ellipsis
type EllipsisMatcherB struct {
	_ *ast.Ellipsis

	Ellipsis token.Pos
	Elt      ExprMatcher
}

func (m EllipsisMatcherB) MatchExpr(cx *MatchContext, node ast.Expr) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m EllipsisMatcherB) Match(cx *MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.Ellipsis)
	if !ok {
		return false, nil
	}
	ok, err = match(cx, ok, err, m.Elt, x.Elt)
	return ok, err
}

// FuncLit
type FuncLitMatcherB struct {
	_ *ast.FuncLit

	Type FuncTypeMatcher
	Body BlockStmtMatcher
}

func (m FuncLitMatcherB) MatchExpr(cx *MatchContext, node ast.Expr) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m FuncLitMatcherB) Match(cx *MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.FuncLit)
	if !ok {
		return false, nil
	}
	ok, err = match(cx, ok, err, m.Type, x.Type)
	ok, err = match(cx, ok, err, m.Body, x.Body)
	return ok, err
}

// FuncType
type FuncTypeMatcherB struct {
	_ *ast.FuncType

	Func       token.Pos
	TypeParams FieldListMatcher
	Params     FieldListMatcher
	Results    FieldListMatcher
}

func (m FuncTypeMatcherB) MatchExpr(cx *MatchContext, node ast.Expr) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m FuncTypeMatcherB) Match(cx *MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.FuncType)
	if !ok {
		return false, nil
	}
	ok, err = match(cx, ok, err, m.TypeParams, x.TypeParams)
	ok, err = match(cx, ok, err, m.Params, x.Params)
	ok, err = match(cx, ok, err, m.Results, x.Results)
	return ok, err
}

// Ident
type IdentMatcherB struct {
	_ *ast.Ident

	NamePos token.Pos
	Name    StringMatcher
}

func (m IdentMatcherB) MatchExpr(cx *MatchContext, node ast.Expr) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m IdentMatcherB) Match(cx *MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.Ident)
	if !ok {
		return false, nil
	}
	ok, err = matchValue(cx, ok, err, m.Name, x.Name)
	return ok, err
}

// IndexExpr
type IndexExprMatcherB struct {
	_ *ast.IndexExpr

	X      ExprMatcher
	Lbrack token.Pos
	Index  ExprMatcher
	Rbrack token.Pos
}

func (m IndexExprMatcherB) MatchExpr(cx *MatchContext, node ast.Expr) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m IndexExprMatcherB) Match(cx *MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.IndexExpr)
	if !ok {
		return false, nil
	}
	ok, err = match(cx, ok, err, m.X, x.X)
	ok, err = match(cx, ok, err, m.Index, x.Index)
	return ok, err
}

// IndexListExpr
type IndexListExprMatcherB struct {
	_ *ast.IndexListExpr

	X       ExprMatcher
	Lbrack  token.Pos
	Indices ListMatcher[ast.Expr]
	Rbrack  token.Pos
}

func (m IndexListExprMatcherB) MatchExpr(cx *MatchContext, node ast.Expr) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m IndexListExprMatcherB) Match(cx *MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.IndexListExpr)
	if !ok {
		return false, nil
	}
	ok, err = match(cx, ok, err, m.X, x.X)
	ok, err = matchList(cx, ok, err, m.Indices, x.Indices)
	return ok, err
}

// InterfaceType
type InterfaceTypeMatcherB struct {
	_ *ast.InterfaceType

	Interface  token.Pos
	Methods    FieldListMatcher
	Incomplete BoolMatcher
}

func (m InterfaceTypeMatcherB) MatchExpr(cx *MatchContext, node ast.Expr) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m InterfaceTypeMatcherB) Match(cx *MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.InterfaceType)
	if !ok {
		return false, nil
	}
	ok, err = match(cx, ok, err, m.Methods, x.Methods)
	ok, err = matchValue(cx, ok, err, m.Incomplete, x.Incomplete)
	return ok, err
}

// KeyValueExpr
type KeyValueExprMatcherB struct {
	_ *ast.KeyValueExpr

	Key   ExprMatcher
	Colon token.Pos
	Value ExprMatcher
}

func (m KeyValueExprMatcherB) MatchExpr(cx *MatchContext, node ast.Expr) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m KeyValueExprMatcherB) Match(cx *MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.KeyValueExpr)
	if !ok {
		return false, nil
	}
	ok, err = match(cx, ok, err, m.Key, x.Key)
	ok, err = match(cx, ok, err, m.Value, x.Value)
	return ok, err
}

// MapType
type MapTypeMatcherB struct {
	_ *ast.MapType

	Map   token.Pos
	Key   ExprMatcher
	Value ExprMatcher
}

func (m MapTypeMatcherB) MatchExpr(cx *MatchContext, node ast.Expr) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m MapTypeMatcherB) Match(cx *MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.MapType)
	if !ok {
		return false, nil
	}
	ok, err = match(cx, ok, err, m.Key, x.Key)
	ok, err = match(cx, ok, err, m.Value, x.Value)
	return ok, err
}

// ParenExpr
type ParenExprMatcherB struct {
	_ *ast.ParenExpr

	Lparen token.Pos
	X      ExprMatcher
	Rparen token.Pos
}

func (m ParenExprMatcherB) MatchExpr(cx *MatchContext, node ast.Expr) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m ParenExprMatcherB) Match(cx *MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.ParenExpr)
	if !ok {
		return false, nil
	}
	ok, err = match(cx, ok, err, m.X, x.X)
	return ok, err
}

// SelectorExpr
type SelectorExprMatcherB struct {
	_ *ast.SelectorExpr

	X   ExprMatcher
	Sel IdentMatcher
}

func (m SelectorExprMatcherB) MatchExpr(cx *MatchContext, node ast.Expr) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m SelectorExprMatcherB) Match(cx *MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.SelectorExpr)
	if !ok {
		return false, nil
	}
	ok, err = match(cx, ok, err, m.X, x.X)
	ok, err = match(cx, ok, err, m.Sel, x.Sel)
	return ok, err
}

// SliceExpr
type SliceExprMatcherB struct {
	_ *ast.SliceExpr

	X      ExprMatcher
	Lbrack token.Pos
	Low    ExprMatcher
	High   ExprMatcher
	Max    ExprMatcher
	Slice3 BoolMatcher
	Rbrack token.Pos
}

func (m SliceExprMatcherB) MatchExpr(cx *MatchContext, node ast.Expr) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m SliceExprMatcherB) Match(cx *MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.SliceExpr)
	if !ok {
		return false, nil
	}
	ok, err = match(cx, ok, err, m.X, x.X)
	ok, err = match(cx, ok, err, m.Low, x.Low)
	ok, err = match(cx, ok, err, m.High, x.High)
	ok, err = match(cx, ok, err, m.Max, x.Max)
	ok, err = matchValue(cx, ok, err, m.Slice3, x.Slice3)
	return ok, err
}

// StarExpr
type StarExprMatcherB struct {
	_ *ast.StarExpr

	Star token.Pos
	X    ExprMatcher
}

func (m StarExprMatcherB) MatchExpr(cx *MatchContext, node ast.Expr) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m StarExprMatcherB) Match(cx *MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.StarExpr)
	if !ok {
		return false, nil
	}
	ok, err = match(cx, ok, err, m.X, x.X)
	return ok, err
}

// StructType
type StructTypeMatcherB struct {
	_ *ast.StructType

	Struct     token.Pos
	Fields     FieldListMatcher
	Incomplete BoolMatcher
}

func (m StructTypeMatcherB) MatchExpr(cx *MatchContext, node ast.Expr) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m StructTypeMatcherB) Match(cx *MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.StructType)
	if !ok {
		return false, nil
	}
	ok, err = match(cx, ok, err, m.Fields, x.Fields)
	ok, err = matchValue(cx, ok, err, m.Incomplete, x.Incomplete)
	return ok, err
}

// TypeAssertExpr
type TypeAssertExprMatcherB struct {
	_ *ast.TypeAssertExpr

	X      ExprMatcher
	Lparen token.Pos
	Type   ExprMatcher
	Rparen token.Pos
}

func (m TypeAssertExprMatcherB) MatchExpr(cx *MatchContext, node ast.Expr) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m TypeAssertExprMatcherB) Match(cx *MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.TypeAssertExpr)
	if !ok {
		return false, nil
	}
	ok, err = match(cx, ok, err, m.X, x.X)
	ok, err = match(cx, ok, err, m.Type, x.Type)
	return ok, err
}

// UnaryExpr
type UnaryExprMatcherB struct {
	_ *ast.UnaryExpr

	OpPos token.Pos
	Op    token.Token
	X     ExprMatcher
}

func (m UnaryExprMatcherB) MatchExpr(cx *MatchContext, node ast.Expr) (ok bool, err error) {
	return m.Match(cx, node)
}
func (m UnaryExprMatcherB) Match(cx *MatchContext, node ast.Node) (ok bool, err error) {
	x, ok := node.(*ast.UnaryExpr)
	if !ok {
		return false, nil
	}
	ok, err = match(cx, ok, err, m.X, x.X)
	return ok, err
}
