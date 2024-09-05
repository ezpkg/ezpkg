package codez

import (
	"go/ast"
)

type VisitContext struct {
}

func Walk(node ast.Node, fn func(cx *VisitContext, node ast.Node) bool) {
	v := zVisitor{
		fn: fn,
	}
	v.visit(node)
}

type zVisitFunc func(cx *VisitContext, node ast.Node) bool

type zVisitor struct {
	cx *VisitContext
	fn zVisitFunc
}

func (v *zVisitor) visit(node ast.Node) {
	switch x := node.(type) {
	case ast.Expr:
		v.visitExpr(x)
	case ast.Stmt:
		v.visitStmt(x)
	case ast.Decl:
		v.visitDecl(x)
	case ast.Spec:
		v.visitSpec(x)
	case ast.Node:
		v.visitOther(x)
	}
}

func (v *zVisitor) visitFile(node *ast.File) {
	ok := v.fn(v.cx, node)
	if !ok {
		return
	}
	v.visitCommentGroup(node.Doc)
	v.visitIdent(node.Name)
	for _, decl := range node.Decls {
		v.visitDecl(decl)
	}
	for _, imp := range node.Imports {
		v.visitSpec(imp)
	}
	for _, cmt := range node.Comments {
		v.visitCommentGroup(cmt)
	}
}

func (v *zVisitor) visitCommentGroup(node *ast.CommentGroup) {
	if node == nil {
		return
	}
	ok := v.fn(v.cx, node)
	if !ok {
		return
	}
	for _, cmt := range node.List {
		v.visitComment(cmt)
	}
}

func (v *zVisitor) visitComment(node *ast.Comment) {
	v.fn(v.cx, node)
}

func (v *zVisitor) visitFieldList(node *ast.FieldList) {
	if node == nil {
		return
	}
	ok := v.fn(v.cx, node)
	if !ok {
		return
	}
	for _, field := range node.List {
		v.visitField(field)
	}
}

func (v *zVisitor) visitField(node *ast.Field) {
	ok := v.fn(v.cx, node)
	if !ok {
		return
	}
	v.visitCommentGroup(node.Doc)
	for _, name := range node.Names {
		v.visitIdent(name)
	}
	v.visitExpr(node.Type)
	v.visitBasicLit(node.Tag)
	v.visitCommentGroup(node.Comment)
}
