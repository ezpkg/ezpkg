//go:build !genz

// Code generated by genz codez-matchers. DO NOT EDIT.

package codez

import (
	ast "go/ast"
)

func (v *zVisitor) visitDecl(node ast.Decl) {
	switch x := node.(type) {
	case *ast.FuncDecl:
		v.visitFuncDecl(x)
	case *ast.GenDecl:
		v.visitGenDecl(x)
	default:
		panic("unreachable ❌")
	}
}
func (v *zVisitor) visitExpr(node ast.Expr) {
	switch x := node.(type) {
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
		panic("unreachable ❌")
	}
}
func (v *zVisitor) visitSpec(node ast.Spec) {
	switch x := node.(type) {
	case *ast.ImportSpec:
		v.visitImportSpec(x)
	case *ast.TypeSpec:
		v.visitTypeSpec(x)
	case *ast.ValueSpec:
		v.visitValueSpec(x)
	default:
		panic("unreachable ❌")
	}
}
func (v *zVisitor) visitStmt(node ast.Stmt) {
	switch x := node.(type) {
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
		panic("unreachable ❌")
	}
}
func (v *zVisitor) visitOther(node ast.Node) {
	switch x := node.(type) {
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
	case *ast.ImportSpec:
		v.visitImportSpec(x)
	case *ast.TypeSpec:
		v.visitTypeSpec(x)
	case *ast.ValueSpec:
		v.visitValueSpec(x)
	default:
		panic("unreachable ❌")
	}
}

// ArrayType
func (v *zVisitor) visitArrayType(node *ast.ArrayType) {
	v.fn(v.cx, node.Len)
	v.visitExpr(node.Len)
	v.fn(v.cx, node.Elt)
	v.visitExpr(node.Elt)
}

// AssignStmt
func (v *zVisitor) visitAssignStmt(node *ast.AssignStmt) {
	for _, item := range node.Lhs {
		ok := v.fn(v.cx, item)
		if ok {
			v.visitExpr(item)
		}
	}
	for _, item := range node.Rhs {
		ok := v.fn(v.cx, item)
		if ok {
			v.visitExpr(item)
		}
	}
}

// BasicLit
func (v *zVisitor) visitBasicLit(node *ast.BasicLit) {
}

// BinaryExpr
func (v *zVisitor) visitBinaryExpr(node *ast.BinaryExpr) {
	v.fn(v.cx, node.X)
	v.visitExpr(node.X)
	v.fn(v.cx, node.Y)
	v.visitExpr(node.Y)
}

// BlockStmt
func (v *zVisitor) visitBlockStmt(node *ast.BlockStmt) {
	for _, item := range node.List {
		ok := v.fn(v.cx, item)
		if ok {
			v.visitStmt(item)
		}
	}
}

// BranchStmt
func (v *zVisitor) visitBranchStmt(node *ast.BranchStmt) {
	v.fn(v.cx, node.Label)
	v.visitIdent(node.Label)
}

// CallExpr
func (v *zVisitor) visitCallExpr(node *ast.CallExpr) {
	v.fn(v.cx, node.Fun)
	v.visitExpr(node.Fun)
	for _, item := range node.Args {
		ok := v.fn(v.cx, item)
		if ok {
			v.visitExpr(item)
		}
	}
}

// CaseClause
func (v *zVisitor) visitCaseClause(node *ast.CaseClause) {
	for _, item := range node.List {
		ok := v.fn(v.cx, item)
		if ok {
			v.visitExpr(item)
		}
	}
	for _, item := range node.Body {
		ok := v.fn(v.cx, item)
		if ok {
			v.visitStmt(item)
		}
	}
}

// ChanType
func (v *zVisitor) visitChanType(node *ast.ChanType) {
	v.fn(v.cx, node.Value)
	v.visitExpr(node.Value)
}

// CommClause
func (v *zVisitor) visitCommClause(node *ast.CommClause) {
	v.fn(v.cx, node.Comm)
	v.visitStmt(node.Comm)
	for _, item := range node.Body {
		ok := v.fn(v.cx, item)
		if ok {
			v.visitStmt(item)
		}
	}
}

// CompositeLit
func (v *zVisitor) visitCompositeLit(node *ast.CompositeLit) {
	v.fn(v.cx, node.Type)
	v.visitExpr(node.Type)
	for _, item := range node.Elts {
		ok := v.fn(v.cx, item)
		if ok {
			v.visitExpr(item)
		}
	}
}

// DeclStmt
func (v *zVisitor) visitDeclStmt(node *ast.DeclStmt) {
	v.fn(v.cx, node.Decl)
	v.visitDecl(node.Decl)
}

// DeferStmt
func (v *zVisitor) visitDeferStmt(node *ast.DeferStmt) {
	v.fn(v.cx, node.Call)
	v.visitCallExpr(node.Call)
}

// Ellipsis
func (v *zVisitor) visitEllipsis(node *ast.Ellipsis) {
	v.fn(v.cx, node.Elt)
	v.visitExpr(node.Elt)
}

// EmptyStmt
func (v *zVisitor) visitEmptyStmt(node *ast.EmptyStmt) {
}

// ExprStmt
func (v *zVisitor) visitExprStmt(node *ast.ExprStmt) {
	v.fn(v.cx, node.X)
	v.visitExpr(node.X)
}

// ForStmt
func (v *zVisitor) visitForStmt(node *ast.ForStmt) {
	v.fn(v.cx, node.Init)
	v.visitStmt(node.Init)
	v.fn(v.cx, node.Cond)
	v.visitExpr(node.Cond)
	v.fn(v.cx, node.Post)
	v.visitStmt(node.Post)
	v.fn(v.cx, node.Body)
	v.visitBlockStmt(node.Body)
}

// FuncDecl
func (v *zVisitor) visitFuncDecl(node *ast.FuncDecl) {
	v.fn(v.cx, node.Doc)
	v.visitCommentGroup(node.Doc)
	v.fn(v.cx, node.Recv)
	v.visitFieldList(node.Recv)
	v.fn(v.cx, node.Name)
	v.visitIdent(node.Name)
	v.fn(v.cx, node.Type)
	v.visitFuncType(node.Type)
	v.fn(v.cx, node.Body)
	v.visitBlockStmt(node.Body)
}

// FuncLit
func (v *zVisitor) visitFuncLit(node *ast.FuncLit) {
	v.fn(v.cx, node.Type)
	v.visitFuncType(node.Type)
	v.fn(v.cx, node.Body)
	v.visitBlockStmt(node.Body)
}

// FuncType
func (v *zVisitor) visitFuncType(node *ast.FuncType) {
	v.fn(v.cx, node.TypeParams)
	v.visitFieldList(node.TypeParams)
	v.fn(v.cx, node.Params)
	v.visitFieldList(node.Params)
	v.fn(v.cx, node.Results)
	v.visitFieldList(node.Results)
}

// GenDecl
func (v *zVisitor) visitGenDecl(node *ast.GenDecl) {
	v.fn(v.cx, node.Doc)
	v.visitCommentGroup(node.Doc)
	for _, item := range node.Specs {
		ok := v.fn(v.cx, item)
		if ok {
			v.visitSpec(item)
		}
	}
}

// GoStmt
func (v *zVisitor) visitGoStmt(node *ast.GoStmt) {
	v.fn(v.cx, node.Call)
	v.visitCallExpr(node.Call)
}

// Ident
func (v *zVisitor) visitIdent(node *ast.Ident) {
}

// IfStmt
func (v *zVisitor) visitIfStmt(node *ast.IfStmt) {
	v.fn(v.cx, node.Init)
	v.visitStmt(node.Init)
	v.fn(v.cx, node.Cond)
	v.visitExpr(node.Cond)
	v.fn(v.cx, node.Body)
	v.visitBlockStmt(node.Body)
	v.fn(v.cx, node.Else)
	v.visitStmt(node.Else)
}

// ImportSpec
func (v *zVisitor) visitImportSpec(node *ast.ImportSpec) {
	v.fn(v.cx, node.Doc)
	v.visitCommentGroup(node.Doc)
	v.fn(v.cx, node.Name)
	v.visitIdent(node.Name)
	v.fn(v.cx, node.Path)
	v.visitBasicLit(node.Path)
	v.fn(v.cx, node.Comment)
	v.visitCommentGroup(node.Comment)
}

// IncDecStmt
func (v *zVisitor) visitIncDecStmt(node *ast.IncDecStmt) {
	v.fn(v.cx, node.X)
	v.visitExpr(node.X)
}

// IndexExpr
func (v *zVisitor) visitIndexExpr(node *ast.IndexExpr) {
	v.fn(v.cx, node.X)
	v.visitExpr(node.X)
	v.fn(v.cx, node.Index)
	v.visitExpr(node.Index)
}

// IndexListExpr
func (v *zVisitor) visitIndexListExpr(node *ast.IndexListExpr) {
	v.fn(v.cx, node.X)
	v.visitExpr(node.X)
	for _, item := range node.Indices {
		ok := v.fn(v.cx, item)
		if ok {
			v.visitExpr(item)
		}
	}
}

// InterfaceType
func (v *zVisitor) visitInterfaceType(node *ast.InterfaceType) {
	v.fn(v.cx, node.Methods)
	v.visitFieldList(node.Methods)
}

// KeyValueExpr
func (v *zVisitor) visitKeyValueExpr(node *ast.KeyValueExpr) {
	v.fn(v.cx, node.Key)
	v.visitExpr(node.Key)
	v.fn(v.cx, node.Value)
	v.visitExpr(node.Value)
}

// LabeledStmt
func (v *zVisitor) visitLabeledStmt(node *ast.LabeledStmt) {
	v.fn(v.cx, node.Label)
	v.visitIdent(node.Label)
	v.fn(v.cx, node.Stmt)
	v.visitStmt(node.Stmt)
}

// MapType
func (v *zVisitor) visitMapType(node *ast.MapType) {
	v.fn(v.cx, node.Key)
	v.visitExpr(node.Key)
	v.fn(v.cx, node.Value)
	v.visitExpr(node.Value)
}

// ParenExpr
func (v *zVisitor) visitParenExpr(node *ast.ParenExpr) {
	v.fn(v.cx, node.X)
	v.visitExpr(node.X)
}

// RangeStmt
func (v *zVisitor) visitRangeStmt(node *ast.RangeStmt) {
	v.fn(v.cx, node.Key)
	v.visitExpr(node.Key)
	v.fn(v.cx, node.Value)
	v.visitExpr(node.Value)
	v.fn(v.cx, node.X)
	v.visitExpr(node.X)
	v.fn(v.cx, node.Body)
	v.visitBlockStmt(node.Body)
}

// ReturnStmt
func (v *zVisitor) visitReturnStmt(node *ast.ReturnStmt) {
	for _, item := range node.Results {
		ok := v.fn(v.cx, item)
		if ok {
			v.visitExpr(item)
		}
	}
}

// SelectStmt
func (v *zVisitor) visitSelectStmt(node *ast.SelectStmt) {
	v.fn(v.cx, node.Body)
	v.visitBlockStmt(node.Body)
}

// SelectorExpr
func (v *zVisitor) visitSelectorExpr(node *ast.SelectorExpr) {
	v.fn(v.cx, node.X)
	v.visitExpr(node.X)
	v.fn(v.cx, node.Sel)
	v.visitIdent(node.Sel)
}

// SendStmt
func (v *zVisitor) visitSendStmt(node *ast.SendStmt) {
	v.fn(v.cx, node.Chan)
	v.visitExpr(node.Chan)
	v.fn(v.cx, node.Value)
	v.visitExpr(node.Value)
}

// SliceExpr
func (v *zVisitor) visitSliceExpr(node *ast.SliceExpr) {
	v.fn(v.cx, node.X)
	v.visitExpr(node.X)
	v.fn(v.cx, node.Low)
	v.visitExpr(node.Low)
	v.fn(v.cx, node.High)
	v.visitExpr(node.High)
	v.fn(v.cx, node.Max)
	v.visitExpr(node.Max)
}

// StarExpr
func (v *zVisitor) visitStarExpr(node *ast.StarExpr) {
	v.fn(v.cx, node.X)
	v.visitExpr(node.X)
}

// StructType
func (v *zVisitor) visitStructType(node *ast.StructType) {
	v.fn(v.cx, node.Fields)
	v.visitFieldList(node.Fields)
}

// SwitchStmt
func (v *zVisitor) visitSwitchStmt(node *ast.SwitchStmt) {
	v.fn(v.cx, node.Init)
	v.visitStmt(node.Init)
	v.fn(v.cx, node.Tag)
	v.visitExpr(node.Tag)
	v.fn(v.cx, node.Body)
	v.visitBlockStmt(node.Body)
}

// TypeAssertExpr
func (v *zVisitor) visitTypeAssertExpr(node *ast.TypeAssertExpr) {
	v.fn(v.cx, node.X)
	v.visitExpr(node.X)
	v.fn(v.cx, node.Type)
	v.visitExpr(node.Type)
}

// TypeSpec
func (v *zVisitor) visitTypeSpec(node *ast.TypeSpec) {
	v.fn(v.cx, node.Doc)
	v.visitCommentGroup(node.Doc)
	v.fn(v.cx, node.Name)
	v.visitIdent(node.Name)
	v.fn(v.cx, node.TypeParams)
	v.visitFieldList(node.TypeParams)
	v.fn(v.cx, node.Type)
	v.visitExpr(node.Type)
	v.fn(v.cx, node.Comment)
	v.visitCommentGroup(node.Comment)
}

// TypeSwitchStmt
func (v *zVisitor) visitTypeSwitchStmt(node *ast.TypeSwitchStmt) {
	v.fn(v.cx, node.Init)
	v.visitStmt(node.Init)
	v.fn(v.cx, node.Assign)
	v.visitStmt(node.Assign)
	v.fn(v.cx, node.Body)
	v.visitBlockStmt(node.Body)
}

// UnaryExpr
func (v *zVisitor) visitUnaryExpr(node *ast.UnaryExpr) {
	v.fn(v.cx, node.X)
	v.visitExpr(node.X)
}

// ValueSpec
func (v *zVisitor) visitValueSpec(node *ast.ValueSpec) {
	v.fn(v.cx, node.Doc)
	v.visitCommentGroup(node.Doc)
	for _, item := range node.Names {
		ok := v.fn(v.cx, item)
		if ok {
			v.visitIdent(item)
		}
	}
	v.fn(v.cx, node.Type)
	v.visitExpr(node.Type)
	for _, item := range node.Values {
		ok := v.fn(v.cx, item)
		if ok {
			v.visitExpr(item)
		}
	}
	v.fn(v.cx, node.Comment)
	v.visitCommentGroup(node.Comment)
}
