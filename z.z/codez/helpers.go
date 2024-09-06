package codez

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"reflect"

	"ezpkg.io/errorz"
)

func NewIdent(name string) *ast.Ident {
	return &ast.Ident{Name: name}
}
func BasicString(value string) *ast.BasicLit {
	return &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("%q", value)}
}
func BasicInt(value int) *ast.BasicLit {
	return &ast.BasicLit{Kind: token.INT, Value: fmt.Sprintf("%d", value)}
}
func BasicFloat(value float64) *ast.BasicLit {
	return &ast.BasicLit{Kind: token.FLOAT, Value: fmt.Sprintf("%v", value)}
}
func BasicBool(value bool) *ast.BasicLit {
	return &ast.BasicLit{Kind: token.IDENT, Value: fmt.Sprintf("%v", value)}
}
func BasicChar(value rune) *ast.BasicLit {
	return &ast.BasicLit{Kind: token.CHAR, Value: fmt.Sprintf("'%c'", value)}
}
func BasicImag(value complex128) *ast.BasicLit {
	return &ast.BasicLit{Kind: token.IMAG, Value: fmt.Sprintf("%v", value)}
}
func NewSelectorExpr(x ast.Expr, sel string) *ast.SelectorExpr {
	return &ast.SelectorExpr{X: x, Sel: NewIdent(sel)}
}

func PrintNode(fset *token.FileSet, node _NodeI) {
	pos := fset.Position(node.Pos())
	if pos.IsValid() {
		fmt.Printf("%v:%v\n", pos.Filename, pos.Line)
	}
	errorz.MustZ(ast.Fprint(os.Stdout, fset, unwrapNode(node), astFilter))
}

func astFilter(name string, value reflect.Value) bool {
	if !ast.NotNilFilter(name, value) {
		return false
	}
	if name == "Pos" || name == "End" {
		return false
	}
	if name == "Obj" || name == "Scope" {
		return false
	}
	return true
}
