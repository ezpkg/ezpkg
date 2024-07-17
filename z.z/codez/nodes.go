package codez

import (
	"go/ast"
)

var otherNodes = []ast.Node{
	&ast.File{},
	&ast.Comment{},
	&ast.CommentGroup{},
}
var allSpecs = []ast.Spec{
	&ast.ImportSpec{},
	&ast.ValueSpec{},
	&ast.TypeSpec{},
}
var allDecls = []ast.Node{
	&ast.BadDecl{},
	&ast.GenDecl{},
	&ast.FuncDecl{},
}
var allStmts = []ast.Stmt{
	&ast.BadStmt{},
	&ast.DeclStmt{},
	&ast.EmptyStmt{},
	&ast.LabeledStmt{},
	&ast.ExprStmt{},
	&ast.SendStmt{},
	&ast.IncDecStmt{},
	&ast.AssignStmt{},
	&ast.GoStmt{},
	&ast.DeferStmt{},
	&ast.ReturnStmt{},
	&ast.BranchStmt{},
	&ast.BlockStmt{},
	&ast.IfStmt{},
	&ast.CaseClause{},
	&ast.SwitchStmt{},
	&ast.TypeSwitchStmt{},
	&ast.CommClause{},
	&ast.SelectStmt{},
	&ast.ForStmt{},
	&ast.RangeStmt{},
}
var allExprs = []ast.Expr{
	&ast.BadExpr{},
	&ast.Ident{},
	&ast.Ellipsis{},
	&ast.BasicLit{},
	&ast.FuncLit{},
	&ast.CompositeLit{},
	&ast.ParenExpr{},
	&ast.SelectorExpr{},
	&ast.IndexExpr{},
	&ast.IndexListExpr{},
	&ast.SliceExpr{},
	&ast.TypeAssertExpr{},
	&ast.CallExpr{},
	&ast.StarExpr{},
	&ast.UnaryExpr{},
	&ast.BinaryExpr{},
	&ast.KeyValueExpr{},
	&ast.ArrayType{},
	&ast.StructType{},
	&ast.FuncType{},
	&ast.InterfaceType{},
	&ast.MapType{},
	&ast.ChanType{},
}
