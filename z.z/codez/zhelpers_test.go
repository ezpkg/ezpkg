package codez

import (
	"strings"

	"ezpkg.io/errorz"
)

var cachePkgs = map[string]*Packages{} // map[pattern]packages

func cacheLoadPackages(pattern ...string) *Packages {
	key := strings.Join(pattern, "|")
	if pkgs := cachePkgs[key]; pkgs != nil {
		return pkgs
	}
	pkgs := errorz.Must(LoadPackages(pattern...))
	return pkgs
}
