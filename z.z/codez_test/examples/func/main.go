package main

import (
	"fmt"

	"ezpkg.io/codez"
	"ezpkg.io/errorz"
)

func main() {
	pkgs := errorz.Must(codez.LoadPackages("ezpkg/z.z/codez_test/testdata/logging/main"))
	for _, pkg := range pkgs {
		fmt.Printf("loaded package %q\n", pkg.PkgPath)
	}
}
