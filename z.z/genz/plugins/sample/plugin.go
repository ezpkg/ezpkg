package sample

import (
	"ezpkg.io/genz"
)

func New() genz.Plugin {
	return plugin{}
}

var _ genz.Filterer = plugin{}

type plugin struct{}

func (p plugin) Name() string { return "sample" }

func (p plugin) Filter(ft genz.FilterEngine) error {
	for _, pkg := range ft.ParsingPackages() {
		ft.IncludePackage(pkg.PkgPath)
		ft.Debugw("include package", "pkg", pkg.PkgPath)
	}
	return nil
}

func (p plugin) Generate(ng genz.Engine) error {
	pkgs := ng.GeneratingPackages()
	for _, gpkg := range pkgs {
		ng.Infow("generate package", "pkg", gpkg.Package.PkgPath)
		objects := gpkg.GetObjects()
		for _, obj := range objects {
			ng.Debugw("  object", "name", obj.Name(), "type", obj.Type())
		}
	}
	return nil
}
