package main

import (
	"fmt"

	"ezpkg.io/codez"
	"ezpkg.io/colorz"
	"ezpkg.io/errorz"
)

func main() {
	pkgs := errorz.Must(codez.LoadPackages("ezpkg.io/-/codez_test/testdata/logging/main"))
	fmt.Println(colorz.Blue.Wrap("👉 loaded packages:"))
	for _, pkg := range pkgs.Packages() {
		fmt.Printf("\t%v\n", pkg.PkgPath)
	}
	fmt.Println()
	fmt.Println(colorz.Blue.Wrap("👉 all packages:"))
	for _, pkg := range pkgs.AllPackages() {
		fmt.Printf("\t%v\n", pkg.PkgPath)
	}
	fmt.Println()
	fmt.Println(colorz.Blue.Wrap("👉 filter ezpkg.io/... , golang.org/..."))
	for _, pkg := range pkgs.AllPackages("ezpkg.io/...", "golang.org/...") {
		fmt.Printf("\t%v\n", pkg.PkgPath)
	}

	matchContext(pkgs)
}

func matchContext(pkgs *codez.Packages) {
	m0 := &codez.SelectorExprMatcher{
		X: &codez.IdentMatcher{
			Name: "stdctx",
		},
		Sel: &codez.IdentMatcher{
			Name: "Context",
		},
	}
	m1 := &codez.IdentMatcher{Name: "ctx"}
	_, _ = m0, m1

	_pkgs := pkgs.AllPackages("ezpkg.io/-/codez_test/testdata/logging/...")
	for _, pkg := range _pkgs {
		nodes := codez.Match(m0, pkg)
		fmt.Printf("\n👉 %v: found %v nodes\n", pkg.PkgPath, len(nodes))
		for _, node := range nodes {
			PrintAst(pkg.Fset, node)
		}
	}
}

func searchError(pkgs *codez.Packages) {
	sr := codez.NewSearch("error").
		InPackages("ezpkg.io/-/codez_test/testdata/logging/...")
	result, err := sr.Exec(pkgs)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	fmt.Println(result)
}

func searchContext(pkgs *codez.Packages) {
	sr := codez.NewSearch("context.Context").
		Import("context", "context").
		InPackages("ezpkg.io/-/codez_test/testdata/logging/...")
	result, err := sr.Exec(pkgs)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	fmt.Println(result)
}

// search for all functions that return error as the last result
func searchFuncReturningError(pkgs *codez.Packages) {
	sr := codez.NewSearch("func $foo(...) (..., error)").
		WithIdent("$foo").
		InPackages("ezpkg.io/-/codez_test/testdata/logging/...")
	result, err := sr.Exec(pkgs)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	fmt.Println(result)
}
