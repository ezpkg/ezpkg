package codez

import (
	"testing"

	g "github.com/onsi/gomega"

	. "ezpkg.io/conveyz"
	"ezpkg.io/diffz"
	"ezpkg.io/errorz"
	"ezpkg.io/stringz"
	"ezpkg.io/testingz"
)

func TestLoadPackages(t *testing.T) {
	ΩxNoDiff := testingz.ConveyDiffByLine(diffz.IgnoreSpace())

	pkgList := func(pkgs []*Package) string {
		return stringz.JoinFunc(pkgs, "\n",
			func(pkg *Package) string { return pkg.PkgPath })
	}
	Convey("LoadPackages", t, func() {
		Convey("single, absolute path", func() {
			pkgs := errorz.Must(LoadPackages("ezpkg.io/-/codez_test/testpkg/logging/main"))

			expected := `ezpkg.io/-/codez_test/testpkg/logging/main`
			ΩxNoDiff(expected, pkgList(pkgs.InputPackages()))
		})
		Convey("absolute path with ...", func() {
			pkgs := errorz.Must(LoadPackages("ezpkg.io/-/codez_test/testpkg/logging/..."))

			expected := `
ezpkg.io/-/codez_test/testpkg/logging
ezpkg.io/-/codez_test/testpkg/logging/main`
			ΩxNoDiff(expected, pkgList(pkgs.InputPackages()))
		})
		Convey("relative path", func() {
			pkgs := errorz.Must(LoadPackages("../codez_test/testpkg/logging/..."))

			expected := `
ezpkg.io/-/codez_test/testpkg/logging
ezpkg.io/-/codez_test/testpkg/logging/main`
			ΩxNoDiff(expected, pkgList(pkgs.InputPackages()))
		})

		pkgs := errorz.Must(LoadPackages("ezpkg.io/-/codez_test/testpkg/...", "golang.org/..."))
		Convey("filter", func() {
			Convey("logging/...", func() {
				zpkgs := pkgs.AllPackages("ezpkg.io/-/codez_test/testpkg/logging/...")

				expected := `
ezpkg.io/-/codez_test/testpkg/logging
ezpkg.io/-/codez_test/testpkg/logging/main`
				ΩxNoDiff(expected, pkgList(zpkgs))
			})
			Convey("golang.org/x/net/html/...", func() {
				zpkgs := pkgs.AllPackages("golang.org/x/net/html/...")

				expected := `
golang.org/x/net/html
golang.org/x/net/html/atom
golang.org/x/net/html/charset`
				ΩxNoDiff(expected, pkgList(zpkgs))
			})
		})
	})
}

func TestPackages(t *testing.T) {
	Ω := GomegaExpect
	pkgs := errorz.Must(LoadPackages("ezpkg.io/-/codez_test/testpkg/logging/main"))

	Convey("Packages", t, func() {
		Convey("GetObject", func() {
			Convey("context.Context", func() {
				objContext := pkgs.GetObject("context", "Context")
				Ω(objContext).ToNot(g.BeNil())
				Ω(objContext.Name()).To(g.Equal("Context"))
				Ω(objContext.Pkg().Path()).To(g.Equal("context"))
				Ω(objContext.Type().String()).To(g.Equal("context.Context"))

				Ω(pkgs.GetType("context", "Context")).To(g.Equal(objContext.Type()))
			})
			Convey("builtin: error", func() {
				objError := pkgs.GetObject("", "error")
				Ω(objError).ToNot(g.BeNil())
				Ω(objError.Name()).To(g.Equal("error"))

				typError := pkgs.GetBuiltInType("error")
				Ω(typError).ToNot(g.BeNil())
				Ω(typError).To(g.Equal(objError.Type()))

				Ω(pkgs.GetType("", "error")).To(g.Equal(typError))
			})
		})

		Convey("GetPackageByPos", func() {
			pkgContext := pkgs.GetPackageByPath("context")
			objContext := pkgs.GetObject("context", "Context")

			Ω(pkgContext).ToNot(g.BeNil())
			Ω(objContext).ToNot(g.BeNil())

			Convey("happy", func() {
				pkg := pkgs.GetPackageByPos(objContext.Pos())
				Ω(pkg).ToNot(g.BeNil())
				Ω(pkg.PkgPath).To(g.Equal("context"))
				Ω(pkg).To(g.Equal(pkgContext))
			})
		})
	})
}
