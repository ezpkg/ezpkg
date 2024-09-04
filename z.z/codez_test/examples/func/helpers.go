package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"reflect"

	"ezpkg.io/errorz"
)

func PrintNode(fset *token.FileSet, node ast.Node) {
	pos := fset.Position(node.Pos())
	if pos.IsValid() {
		fmt.Printf("-> %v:%v\n", pos.Filename, pos.Line)
	}
	errorz.MustZ(ast.Fprint(os.Stdout, fset, node, astFilter))
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
