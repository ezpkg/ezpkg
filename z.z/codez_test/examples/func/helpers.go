package main

import (
	"go/ast"
	"go/token"

	"ezpkg.io/codez"
)

func PrintNode(fset *token.FileSet, node ast.Node) {
	codez.PrintNode(fset, node)
}
