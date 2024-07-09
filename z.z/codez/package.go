package codez

import (
	"slices"
	"strings"

	"golang.org/x/tools/go/packages"

	"ezpkg.io/slicez"
)

type Packages struct {
	pkgs []*Package

	mapPkgs map[string]*Package
	allPkgs []*Package
	stdPkgs []*Package
}

type Package struct {
	*packages.Package
}

func newPackage(pkg *packages.Package) *Package {
	return &Package{
		Package: pkg,
	}
}

func (p *Packages) Packages() []*Package {
	return p.pkgs
}

// AllPackages returns all packages, including std packages, golang.org/x packages, and other packages. It supports filtering by pattern. Examples:
//
//	AllPackages()                         👉 return all packages
//	AllPackages("ezpkg.io/...")           👉 return packages that start with "ezpkg.io"
//	AllPackages("ezpkg.io/codez", "fmt")  👉 return listed packages
func (p *Packages) AllPackages(pattern ...string) []*Package {
	return filterPackages(p.allPkgs, pattern...)
}
func (p *Packages) StdPackages() []*Package {
	return p.stdPkgs
}
func (p *Packages) NonStdPackages() []*Package {
	return p.allPkgs[len(p.stdPkgs):]
}

func newPackages(pkgs []*Package) *Packages {
	isStd := func(path string) bool {
		return !strings.Contains(path, ".")
	}
	isGolangOrg := func(path string) bool {
		return strings.HasPrefix(path, "golang.org/")
	}
	sortPkgs := func(pkgs []*Package) {
		slices.SortFunc(pkgs, func(a, b *Package) int {
			return strings.Compare(a.PkgPath, b.PkgPath)
		})
	}

	p := &Packages{pkgs: pkgs}
	p.mapPkgs = map[string]*Package{}
	for _, pkg := range p.pkgs {
		for path, impPkg := range pkg.Imports {
			if p.mapPkgs[path] == nil {
				pkg0 := newPackage(impPkg)
				p.mapPkgs[path] = pkg0
				p.allPkgs = append(p.allPkgs, pkg0)
			}
		}
	}

	// filter std packages then sort
	var goOrgPkgs, otherPkgs []*Package
	for _, pkg := range p.mapPkgs {
		switch {
		case isStd(pkg.PkgPath):
			p.stdPkgs = append(p.stdPkgs, pkg)
		case isGolangOrg(pkg.PkgPath):
			goOrgPkgs = append(goOrgPkgs, pkg)
		default:
			otherPkgs = append(otherPkgs, pkg)
		}
	}
	sortPkgs(p.stdPkgs)
	sortPkgs(goOrgPkgs)
	sortPkgs(otherPkgs)
	p.allPkgs = slicez.Concat(p.stdPkgs, goOrgPkgs, otherPkgs)
	return p
}

func filterPackages(pkgs []*Package, pattern ...string) []*Package {
	if len(pattern) == 0 {
		return pkgs
	}
	var out []*Package
	for _, p := range pattern {
		if !strings.HasSuffix(p, "/...") {
			for _, pkg := range pkgs {
				if pkg.PkgPath == p {
					out = append(out, pkg)
				}
			}
			continue
		}

		p = strings.TrimSuffix(p, "...")
		for _, pkg := range pkgs {
			if pkg.PkgPath == strings.TrimSuffix(p, "/") {
				out = append(out, pkg)
			} else if strings.HasPrefix(pkg.PkgPath, p) {
				out = append(out, pkg)
			}
		}
	}
	return out
}
