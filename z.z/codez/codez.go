// Package codez provides API for working with go code, including parsing, searching, and refactoring code.
package codez // import "ezpkg.io/codez"

import (
	"golang.org/x/tools/go/packages"

	"ezpkg.io/errorz"
)

func LoadPackages(pattern ...string) (_ *Packages, err error) {
	cfg := Cfg()
	pkgs, err := packages.Load(&cfg.Config, pattern...)
	if err != nil {
		return nil, err
	}

	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			// only store the first error message, use .AllErrors() to get all
			err0 := errorz.NoStack().Wrapf(pkg.Errors[0], "package %q", pkg.PkgPath)
			errorz.AppendTo(&err, err0)
		}
	}
	err = errorz.Wrap(err, "failed to load packages (use .Errors() to get all errors)")
	pkgSet := newPackages(pkgs)
	if pkgSet == nil && err == nil {
		err = errorz.New("no packages loaded")
	}
	return pkgSet, err
}
