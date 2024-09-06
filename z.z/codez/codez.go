// Package codez provides API for working with go code, including parsing, searching, and refactoring code.
package codez // import "ezpkg.io/codez"

import (
	"golang.org/x/tools/go/packages"

	"ezpkg.io/errorz"
)

// LoadPackages loads packages by pattern. To see compilation errors, use Packages.AllErrors(), Packages.FirstErrors(), or Package.Errors.
func LoadPackages(pattern ...string) (_ *Packages, err error) {
	cfg := Cfg()
	pkgs, err := packages.Load(&cfg.Config, pattern...)
	if err != nil {
		return nil, err
	}
	if len(pkgs) == 0 {
		return nil, errorz.New("no packages loaded")
	}
	return newPackages(pkgs)
}
