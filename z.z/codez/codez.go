// Package codez provides API for working with go code, including parsing, searching, and refactoring code.
package codez // import "ezpkg.io/codez"

import (
	"golang.org/x/tools/go/packages"

	"ezpkg.io/errorz"
	"ezpkg.io/slicez"
)

func LoadPackages(pattern ...string) (_ *Packages, err error) {
	cfg := Cfg()
	pkgs, err := packages.Load(&cfg.Config, pattern...)
	if err != nil {
		return nil, err
	}
	zPkgs := slicez.MapFunc(pkgs, func(pkg *packages.Package) *Package {
		if len(pkg.Errors) > 0 {
			// only store the first error message, use .AllErrors() to get all
			errorz.NoStack().AppendTo(&err, pkg.Errors[0])
		}
		return &Package{Package: pkg}
	})
	err = errorz.Wrap(err, "failed to load packages (use .Errors() to get all errors)")
	pkgSet := newPackages(zPkgs)
	if pkgSet == nil && err == nil {
		err = errorz.New("no packages loaded")
	}
	return pkgSet, err
}
