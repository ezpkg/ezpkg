// Package codez provides API for working with go code, including parsing, searching, and refactoring code.
package codez // import "ezpkg.io/codez"

import (
	"golang.org/x/tools/go/packages"

	"ezpkg.io/slicez"
)

type Config struct {
	packages.Config
}

func Cfg() *Config {
	return &Config{
		Config: packages.Config{
			Mode: 0 |
				packages.NeedName |
				packages.NeedFiles |
				packages.NeedCompiledGoFiles |
				packages.NeedImports |
				packages.NeedDeps |
				packages.NeedExportFile |
				packages.NeedTypes |
				packages.NeedSyntax |
				packages.NeedTypesInfo |
				packages.NeedTypesSizes |
				packages.NeedModule |
				packages.NeedEmbedFiles |
				packages.NeedEmbedPatterns,
		},
	}
}

func (c *Config) Unwrap() *packages.Config {
	return &c.Config
}

func LoadPackages(pattern ...string) (Packages, error) {
	cfg := Cfg()
	pkgs, err := packages.Load(&cfg.Config, pattern...)
	if err != nil {
		return nil, err
	}
	return slicez.MapFunc(pkgs, func(pkg *packages.Package) *Package {
		return &Package{Package: pkg}
	}), nil
}
