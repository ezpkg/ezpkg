package main

import (
	"go/ast"
	"go/token"
	"os"
	"reflect"

	"ezpkg.io/errorz"
)

func PrintAst(fset *token.FileSet, x any) {
	errorz.MustZ(ast.Fprint(os.Stdout, fset, x, astFilter))
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
