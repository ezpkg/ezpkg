package codez

import (
	"testing"

	g "github.com/onsi/gomega"

	. "ezpkg.io/conveyz"
	"ezpkg.io/errorz"
)

func TestPackage(t *testing.T) {
	Ω := GomegaExpect
	pkgs := errorz.Must(LoadPackages("ezpkg.io/-/codez_test/testdata/logging/main"))

	Convey("Packages", t, func() {
		Convey("LoadPackages", func() {
			Ω(pkgs.Packages()).To(g.HaveLen(1))

			pkg := pkgs.Packages()[0]
			Ω(pkg.PkgPath).To(g.Equal("ezpkg.io/-/codez_test/testdata/logging/main"))
		})

		Convey("GetObject", func() {
			objContext := pkgs.GetObject("context", "Context")
			Ω(objContext).ToNot(g.BeNil())
			Ω(objContext.Name()).To(g.Equal("Context"))
			Ω(objContext.Pkg().Path()).To(g.Equal("context"))
			Ω(objContext.Type().String()).To(g.Equal("context.Context"))
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
			Convey("fast path", func() {
				pkg0 := quickSearchPkgByPos(pkgs.pkgByPos, objContext.Pos())
				Ω(pkg0).To(g.Equal(pkgContext))
			})
		})
	})
}
