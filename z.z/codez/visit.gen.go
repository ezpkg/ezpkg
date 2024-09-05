//go:build !genz

// Code generated by genz codez-matchers. DO NOT EDIT.

package codez

import (
	"fmt"
	ast "go/ast"
)

func (v *zVisitor) visitDecl(node ast.Decl) {
	switch x := node.(type) {
	case nil:
		return
	case *ast.FuncDecl:
		v.visitFuncDecl(x)
	case *ast.GenDecl:
		v.visitGenDecl(x)
	default:
		panic(fmt.Sprintf("unreachable: %v is ast.Decl ❌", node))
	}
}
func (v *zVisitor) visitExpr(node ast.Expr) {
	switch x := node.(type) {
	case nil:
		return
	case *ast.ArrayType:
		v.visitArrayType(x)
	case *ast.BasicLit:
		v.visitBasicLit(x)
	case *ast.BinaryExpr:
		v.visitBinaryExpr(x)
	case *ast.CallExpr:
		v.visitCallExpr(x)
	case *ast.ChanType:
		v.visitChanType(x)
	case *ast.CompositeLit:
		v.visitCompositeLit(x)
	case *ast.Ellipsis:
		v.visitEllipsis(x)
	case *ast.FuncLit:
		v.visitFuncLit(x)
	case *ast.FuncType:
		v.visitFuncType(x)
	case *ast.Ident:
		v.visitIdent(x)
	case *ast.IndexExpr:
		v.visitIndexExpr(x)
	case *ast.IndexListExpr:
		v.visitIndexListExpr(x)
	case *ast.InterfaceType:
		v.visitInterfaceType(x)
	case *ast.KeyValueExpr:
		v.visitKeyValueExpr(x)
	case *ast.MapType:
		v.visitMapType(x)
	case *ast.ParenExpr:
		v.visitParenExpr(x)
	case *ast.SelectorExpr:
		v.visitSelectorExpr(x)
	case *ast.SliceExpr:
		v.visitSliceExpr(x)
	case *ast.StarExpr:
		v.visitStarExpr(x)
	case *ast.StructType:
		v.visitStructType(x)
	case *ast.TypeAssertExpr:
		v.visitTypeAssertExpr(x)
	case *ast.UnaryExpr:
		v.visitUnaryExpr(x)
	default:
		panic(fmt.Sprintf("unreachable: %v is ast.Expr ❌", node))
	}
}
func (v *zVisitor) visitSpec(node ast.Spec) {
	switch x := node.(type) {
	case nil:
		return
	case *ast.ImportSpec:
		v.visitImportSpec(x)
	case *ast.TypeSpec:
		v.visitTypeSpec(x)
	case *ast.ValueSpec:
		v.visitValueSpec(x)
	default:
		panic(fmt.Sprintf("unreachable: %v is ast.Spec ❌", node))
	}
}
func (v *zVisitor) visitStmt(node ast.Stmt) {
	switch x := node.(type) {
	case nil:
		return
	case *ast.AssignStmt:
		v.visitAssignStmt(x)
	case *ast.BlockStmt:
		v.visitBlockStmt(x)
	case *ast.BranchStmt:
		v.visitBranchStmt(x)
	case *ast.CaseClause:
		v.visitCaseClause(x)
	case *ast.CommClause:
		v.visitCommClause(x)
	case *ast.DeclStmt:
		v.visitDeclStmt(x)
	case *ast.DeferStmt:
		v.visitDeferStmt(x)
	case *ast.EmptyStmt:
		v.visitEmptyStmt(x)
	case *ast.ExprStmt:
		v.visitExprStmt(x)
	case *ast.ForStmt:
		v.visitForStmt(x)
	case *ast.GoStmt:
		v.visitGoStmt(x)
	case *ast.IfStmt:
		v.visitIfStmt(x)
	case *ast.IncDecStmt:
		v.visitIncDecStmt(x)
	case *ast.LabeledStmt:
		v.visitLabeledStmt(x)
	case *ast.RangeStmt:
		v.visitRangeStmt(x)
	case *ast.ReturnStmt:
		v.visitReturnStmt(x)
	case *ast.SelectStmt:
		v.visitSelectStmt(x)
	case *ast.SendStmt:
		v.visitSendStmt(x)
	case *ast.SwitchStmt:
		v.visitSwitchStmt(x)
	case *ast.TypeSwitchStmt:
		v.visitTypeSwitchStmt(x)
	default:
		panic(fmt.Sprintf("unreachable: %v is ast.Stmt ❌", node))
	}
}
func (v *zVisitor) visitOther(node ast.Node) {
	switch x := node.(type) {
	case nil:
		return
	case *ast.Comment:
		v.visitComment(x)
	case *ast.CommentGroup:
		v.visitCommentGroup(x)
	case *ast.Field:
		v.visitField(x)
	case *ast.FieldList:
		v.visitFieldList(x)
	case *ast.File:
		v.visitFile(x)
	default:
		panic(fmt.Sprintf("unreachable: %v is ast.Other ❌", node))
	}
}

// ArrayType
func (v *zVisitor) visitArrayType(node *ast.ArrayType) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitExpr(node.Len)
		v.visitExpr(node.Elt)
	}
}

// AssignStmt
func (v *zVisitor) visitAssignStmt(node *ast.AssignStmt) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		for _, item := range node.Lhs {
			v.visitExpr(item)
		}
		for _, item := range node.Rhs {
			v.visitExpr(item)
		}
	}
}

// BasicLit
func (v *zVisitor) visitBasicLit(node *ast.BasicLit) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
	}
}

// BinaryExpr
func (v *zVisitor) visitBinaryExpr(node *ast.BinaryExpr) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitExpr(node.X)
		v.visitExpr(node.Y)
	}
}

// BlockStmt
func (v *zVisitor) visitBlockStmt(node *ast.BlockStmt) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		for _, item := range node.List {
			v.visitStmt(item)
		}
	}
}

// BranchStmt
func (v *zVisitor) visitBranchStmt(node *ast.BranchStmt) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitIdent(node.Label)
	}
}

// CallExpr
func (v *zVisitor) visitCallExpr(node *ast.CallExpr) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitExpr(node.Fun)
		for _, item := range node.Args {
			v.visitExpr(item)
		}
	}
}

// CaseClause
func (v *zVisitor) visitCaseClause(node *ast.CaseClause) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		for _, item := range node.List {
			v.visitExpr(item)
		}
		for _, item := range node.Body {
			v.visitStmt(item)
		}
	}
}

// ChanType
func (v *zVisitor) visitChanType(node *ast.ChanType) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitExpr(node.Value)
	}
}

// CommClause
func (v *zVisitor) visitCommClause(node *ast.CommClause) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitStmt(node.Comm)
		for _, item := range node.Body {
			v.visitStmt(item)
		}
	}
}

// Comment
func (v *zVisitor) visitComment(node *ast.Comment) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
	}
}

// CommentGroup
func (v *zVisitor) visitCommentGroup(node *ast.CommentGroup) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		for _, item := range node.List {
			v.visitComment(item)
		}
	}
}

// CompositeLit
func (v *zVisitor) visitCompositeLit(node *ast.CompositeLit) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitExpr(node.Type)
		for _, item := range node.Elts {
			v.visitExpr(item)
		}
	}
}

// DeclStmt
func (v *zVisitor) visitDeclStmt(node *ast.DeclStmt) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitDecl(node.Decl)
	}
}

// DeferStmt
func (v *zVisitor) visitDeferStmt(node *ast.DeferStmt) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitCallExpr(node.Call)
	}
}

// Ellipsis
func (v *zVisitor) visitEllipsis(node *ast.Ellipsis) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitExpr(node.Elt)
	}
}

// EmptyStmt
func (v *zVisitor) visitEmptyStmt(node *ast.EmptyStmt) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
	}
}

// ExprStmt
func (v *zVisitor) visitExprStmt(node *ast.ExprStmt) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitExpr(node.X)
	}
}

// Field
func (v *zVisitor) visitField(node *ast.Field) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitCommentGroup(node.Doc)
		for _, item := range node.Names {
			v.visitIdent(item)
		}
		v.visitExpr(node.Type)
		v.visitBasicLit(node.Tag)
		v.visitCommentGroup(node.Comment)
	}
}

// FieldList
func (v *zVisitor) visitFieldList(node *ast.FieldList) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		for _, item := range node.List {
			v.visitField(item)
		}
	}
}

// File
func (v *zVisitor) visitFile(node *ast.File) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitCommentGroup(node.Doc)
		v.visitIdent(node.Name)
		for _, item := range node.Decls {
			v.visitDecl(item)
		}
		for _, item := range node.Imports {
			v.visitImportSpec(item)
		}
		for _, item := range node.Unresolved {
			v.visitIdent(item)
		}
		for _, item := range node.Comments {
			v.visitCommentGroup(item)
		}
	}
}

// ForStmt
func (v *zVisitor) visitForStmt(node *ast.ForStmt) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitStmt(node.Init)
		v.visitExpr(node.Cond)
		v.visitStmt(node.Post)
		v.visitBlockStmt(node.Body)
	}
}

// FuncDecl
func (v *zVisitor) visitFuncDecl(node *ast.FuncDecl) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitCommentGroup(node.Doc)
		v.visitFieldList(node.Recv)
		v.visitIdent(node.Name)
		v.visitFuncType(node.Type)
		v.visitBlockStmt(node.Body)
	}
}

// FuncLit
func (v *zVisitor) visitFuncLit(node *ast.FuncLit) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitFuncType(node.Type)
		v.visitBlockStmt(node.Body)
	}
}

// FuncType
func (v *zVisitor) visitFuncType(node *ast.FuncType) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitFieldList(node.TypeParams)
		v.visitFieldList(node.Params)
		v.visitFieldList(node.Results)
	}
}

// GenDecl
func (v *zVisitor) visitGenDecl(node *ast.GenDecl) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitCommentGroup(node.Doc)
		for _, item := range node.Specs {
			v.visitSpec(item)
		}
	}
}

// GoStmt
func (v *zVisitor) visitGoStmt(node *ast.GoStmt) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitCallExpr(node.Call)
	}
}

// Ident
func (v *zVisitor) visitIdent(node *ast.Ident) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
	}
}

// IfStmt
func (v *zVisitor) visitIfStmt(node *ast.IfStmt) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitStmt(node.Init)
		v.visitExpr(node.Cond)
		v.visitBlockStmt(node.Body)
		v.visitStmt(node.Else)
	}
}

// ImportSpec
func (v *zVisitor) visitImportSpec(node *ast.ImportSpec) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitCommentGroup(node.Doc)
		v.visitIdent(node.Name)
		v.visitBasicLit(node.Path)
		v.visitCommentGroup(node.Comment)
	}
}

// IncDecStmt
func (v *zVisitor) visitIncDecStmt(node *ast.IncDecStmt) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitExpr(node.X)
	}
}

// IndexExpr
func (v *zVisitor) visitIndexExpr(node *ast.IndexExpr) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitExpr(node.X)
		v.visitExpr(node.Index)
	}
}

// IndexListExpr
func (v *zVisitor) visitIndexListExpr(node *ast.IndexListExpr) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitExpr(node.X)
		for _, item := range node.Indices {
			v.visitExpr(item)
		}
	}
}

// InterfaceType
func (v *zVisitor) visitInterfaceType(node *ast.InterfaceType) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitFieldList(node.Methods)
	}
}

// KeyValueExpr
func (v *zVisitor) visitKeyValueExpr(node *ast.KeyValueExpr) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitExpr(node.Key)
		v.visitExpr(node.Value)
	}
}

// LabeledStmt
func (v *zVisitor) visitLabeledStmt(node *ast.LabeledStmt) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitIdent(node.Label)
		v.visitStmt(node.Stmt)
	}
}

// MapType
func (v *zVisitor) visitMapType(node *ast.MapType) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitExpr(node.Key)
		v.visitExpr(node.Value)
	}
}

// ParenExpr
func (v *zVisitor) visitParenExpr(node *ast.ParenExpr) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitExpr(node.X)
	}
}

// RangeStmt
func (v *zVisitor) visitRangeStmt(node *ast.RangeStmt) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitExpr(node.Key)
		v.visitExpr(node.Value)
		v.visitExpr(node.X)
		v.visitBlockStmt(node.Body)
	}
}

// ReturnStmt
func (v *zVisitor) visitReturnStmt(node *ast.ReturnStmt) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		for _, item := range node.Results {
			v.visitExpr(item)
		}
	}
}

// SelectStmt
func (v *zVisitor) visitSelectStmt(node *ast.SelectStmt) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitBlockStmt(node.Body)
	}
}

// SelectorExpr
func (v *zVisitor) visitSelectorExpr(node *ast.SelectorExpr) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitExpr(node.X)
		v.visitIdent(node.Sel)
	}
}

// SendStmt
func (v *zVisitor) visitSendStmt(node *ast.SendStmt) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitExpr(node.Chan)
		v.visitExpr(node.Value)
	}
}

// SliceExpr
func (v *zVisitor) visitSliceExpr(node *ast.SliceExpr) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitExpr(node.X)
		v.visitExpr(node.Low)
		v.visitExpr(node.High)
		v.visitExpr(node.Max)
	}
}

// StarExpr
func (v *zVisitor) visitStarExpr(node *ast.StarExpr) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitExpr(node.X)
	}
}

// StructType
func (v *zVisitor) visitStructType(node *ast.StructType) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitFieldList(node.Fields)
	}
}

// SwitchStmt
func (v *zVisitor) visitSwitchStmt(node *ast.SwitchStmt) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitStmt(node.Init)
		v.visitExpr(node.Tag)
		v.visitBlockStmt(node.Body)
	}
}

// TypeAssertExpr
func (v *zVisitor) visitTypeAssertExpr(node *ast.TypeAssertExpr) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitExpr(node.X)
		v.visitExpr(node.Type)
	}
}

// TypeSpec
func (v *zVisitor) visitTypeSpec(node *ast.TypeSpec) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitCommentGroup(node.Doc)
		v.visitIdent(node.Name)
		v.visitFieldList(node.TypeParams)
		v.visitExpr(node.Type)
		v.visitCommentGroup(node.Comment)
	}
}

// TypeSwitchStmt
func (v *zVisitor) visitTypeSwitchStmt(node *ast.TypeSwitchStmt) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitStmt(node.Init)
		v.visitStmt(node.Assign)
		v.visitBlockStmt(node.Body)
	}
}

// UnaryExpr
func (v *zVisitor) visitUnaryExpr(node *ast.UnaryExpr) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitExpr(node.X)
	}
}

// ValueSpec
func (v *zVisitor) visitValueSpec(node *ast.ValueSpec) {
	ok := node != nil && v.fn(v.cx, node)
	if ok {
		v.visitCommentGroup(node.Doc)
		for _, item := range node.Names {
			v.visitIdent(item)
		}
		v.visitExpr(node.Type)
		for _, item := range node.Values {
			v.visitExpr(item)
		}
		v.visitCommentGroup(node.Comment)
	}
}
