package codez

import (
	"go/ast"
	"strings"
	"testing"

	g "github.com/onsi/gomega"

	. "ezpkg.io/conveyz"
	"ezpkg.io/diffz"
	"ezpkg.io/testingz"
)

func TestVisit(t *testing.T) {
	Ω := GomegaExpect
	ΩxNoDiff := testingz.ConveyDiffByLine(diffz.IgnoreSpace())

	Convey("Walk", t, func() {
		pkgs := cacheLoadPackages("ezpkg.io/-/codez_test/testpkg/logging/...")
		pkgLog := pkgs.MustGetPackageByPath("ezpkg.io/-/codez_test/testpkg/logging")
		objLogger := pkgLog.MustGetObject("Logger")
		astFile := pkgs.MustGetFileByPos(objLogger.Pos())

		// 👉 uncomment to print the ast
		// printAst("logging/log.go", pkgs.Fset, astFile)

		paths := []string{}
		Walk(astFile, func(cx *VisitContext, node ast.Node) bool {
			if ident, ok := node.(*ast.Ident); ok && ident.Name == "TryOne" {
				parent := cx.Parent()
				// printAst("parent of TryOne is FuncDecl", pkgs.Fset, parent)

				funcDecl, isFuncDecl := parent.(*ast.FuncDecl)
				Ω(isFuncDecl).To(g.BeTrue())
				Ω(funcDecl.Name.Name).To(g.Equal("TryOne"))
			}

			paths = append(paths, cx.Path())
			return true
		})
		expected := `
Name
Decls.0
Decls.0.Specs.0
Decls.0.Specs.0.Path
Decls.0.Specs.1
Decls.0.Specs.1.Path
Decls.0.Specs.2
Decls.0.Specs.2.Path
Decls.1
Decls.1.Specs.0
Decls.1.Specs.0.Name
Decls.1.Specs.0.Type
Decls.1.Specs.0.Type.Fields
Decls.1.Specs.0.Type.Fields.List.0
Decls.1.Specs.0.Type.Fields.List.0.Names.0
Decls.1.Specs.0.Type.Fields.List.0.Type
Decls.1.Specs.0.Type.Fields.List.0.Type.X
Decls.1.Specs.0.Type.Fields.List.0.Type.Sel
Decls.2
Decls.2.Name
Decls.2.Type
Decls.2.Type.Params
Decls.2.Type.Params.List.0
Decls.2.Type.Params.List.0.Names.0
Decls.2.Type.Params.List.0.Type
Decls.2.Type.Params.List.0.Type.X
Decls.2.Type.Params.List.0.Type.Sel
Decls.2.Type.Results
Decls.2.Type.Results.List.0
Decls.2.Type.Results.List.0.Type
Decls.2.Type.Results.List.0.Type.X
Decls.2.Type.Results.List.1
Decls.2.Type.Results.List.1.Type
Decls.2.Body
Decls.2.Body.List.0
Decls.2.Body.List.0.Results.0
Decls.2.Body.List.0.Results.0.X
Decls.2.Body.List.0.Results.0.X.Type
Decls.2.Body.List.0.Results.0.X.Elts.0
Decls.2.Body.List.0.Results.0.X.Elts.0.Key
Decls.2.Body.List.0.Results.0.X.Elts.0.Value
Decls.2.Body.List.0.Results.1
Decls.3
Decls.3.Recv
Decls.3.Recv.List.0
Decls.3.Recv.List.0.Names.0
Decls.3.Recv.List.0.Type
Decls.3.Recv.List.0.Type.X
Decls.3.Name
Decls.3.Type
Decls.3.Type.Params
Decls.3.Type.Params.List.0
Decls.3.Type.Params.List.0.Names.0
Decls.3.Type.Params.List.0.Type
Decls.3.Type.Params.List.1
Decls.3.Type.Params.List.1.Names.0
Decls.3.Type.Params.List.1.Type
Decls.3.Type.Params.List.1.Type.Elt
Decls.3.Body
Decls.3.Body.List.0
Decls.3.Body.List.0.X
Decls.3.Body.List.0.X.Fun
Decls.3.Body.List.0.X.Fun.X
Decls.3.Body.List.0.X.Fun.Sel
Decls.3.Body.List.0.X.Args.0
Decls.3.Body.List.0.X.Args.1
Decls.3.Body.List.1
Decls.3.Body.List.1.X
Decls.3.Body.List.1.X.Fun
Decls.3.Body.List.1.X.Fun.X
Decls.3.Body.List.1.X.Fun.Sel
Decls.4
Decls.4.Recv
Decls.4.Recv.List.0
Decls.4.Recv.List.0.Names.0
Decls.4.Recv.List.0.Type
Decls.4.Recv.List.0.Type.X
Decls.4.Name
Decls.4.Type
Decls.4.Type.Params
Decls.4.Type.Params.List.0
Decls.4.Type.Params.List.0.Names.0
Decls.4.Type.Params.List.0.Type
Decls.4.Type.Params.List.0.Type.X
Decls.4.Type.Params.List.0.Type.Sel
Decls.4.Type.Params.List.1
Decls.4.Type.Params.List.1.Names.0
Decls.4.Type.Params.List.1.Type
Decls.4.Type.Params.List.2
Decls.4.Type.Params.List.2.Names.0
Decls.4.Type.Params.List.2.Type
Decls.4.Type.Params.List.2.Type.Elt
Decls.4.Body
Decls.4.Body.List.0
Decls.4.Body.List.0.X
Decls.4.Body.List.0.X.Fun
Decls.4.Body.List.0.X.Fun.X
Decls.4.Body.List.0.X.Fun.Sel
Decls.4.Body.List.0.X.Args.0
Decls.4.Body.List.0.X.Args.1
Decls.4.Body.List.1
Decls.4.Body.List.1.X
Decls.4.Body.List.1.X.Fun
Decls.4.Body.List.1.X.Fun.X
Decls.4.Body.List.1.X.Fun.Sel
Decls.5
Decls.5.Recv
Decls.5.Recv.List.0
Decls.5.Recv.List.0.Names.0
Decls.5.Recv.List.0.Type
Decls.5.Recv.List.0.Type.X
Decls.5.Name
Decls.5.Type
Decls.5.Type.Params
Decls.5.Type.Results
Decls.5.Type.Results.List.0
Decls.5.Type.Results.List.0.Type
Decls.5.Body
Decls.5.Body.List.0
Decls.5.Body.List.0.Results.0
Decls.6
Decls.6.Recv
Decls.6.Recv.List.0
Decls.6.Recv.List.0.Names.0
Decls.6.Recv.List.0.Type
Decls.6.Recv.List.0.Type.X
Decls.6.Name
Decls.6.Type
Decls.6.Type.Params
Decls.6.Type.Results
Decls.6.Type.Results.List.0
Decls.6.Type.Results.List.0.Type
Decls.6.Type.Results.List.0.Type.Elt
Decls.6.Body
Decls.6.Body.List.0
Decls.6.Body.List.0.Results.0
Decls.7
Decls.7.Recv
Decls.7.Recv.List.0
Decls.7.Recv.List.0.Names.0
Decls.7.Recv.List.0.Type
Decls.7.Recv.List.0.Type.X
Decls.7.Name
Decls.7.Type
Decls.7.Type.Params
Decls.7.Type.Results
Decls.7.Type.Results.List.0
Decls.7.Type.Results.List.0.Type
Decls.7.Type.Results.List.1
Decls.7.Type.Results.List.1.Type
Decls.7.Body
Decls.7.Body.List.0
Decls.7.Body.List.0.Results.0
Decls.7.Body.List.0.Results.1
Decls.8
Decls.8.Recv
Decls.8.Recv.List.0
Decls.8.Recv.List.0.Names.0
Decls.8.Recv.List.0.Type
Decls.8.Recv.List.0.Type.X
Decls.8.Name
Decls.8.Type
Decls.8.Type.Params
Decls.8.Type.Results
Decls.8.Type.Results.List.0
Decls.8.Type.Results.List.0.Names.0
Decls.8.Type.Results.List.0.Names.1
Decls.8.Type.Results.List.0.Type
Decls.8.Type.Results.List.1
Decls.8.Type.Results.List.1.Names.0
Decls.8.Type.Results.List.1.Type
Decls.8.Body
Decls.8.Body.List.0
Decls.8.Body.List.0.Results.0
Decls.8.Body.List.0.Results.1
Decls.8.Body.List.0.Results.2
Decls.9
Decls.9.Name
Decls.9.Type
Decls.9.Type.Params
Decls.9.Type.Results
Decls.9.Type.Results.List.0
Decls.9.Type.Results.List.0.Type
Decls.9.Body
Decls.9.Body.List.0
Decls.9.Body.List.0.Results.0
Decls.9.Body.List.0.Results.0.Fun
Decls.9.Body.List.0.Results.0.Fun.X
Decls.9.Body.List.0.Results.0.Fun.Sel
Decls.9.Body.List.0.Results.0.Args.0
Imports.0
Imports.0.Path
Imports.1
Imports.1.Path
Imports.2
Imports.2.Path
`
		ΩxNoDiff(expected, strings.Join(paths, "\n"))
	})
}
