package codez

import (
	"cmp"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"slices"
	"strings"

	"golang.org/x/tools/go/packages"

	"ezpkg.io/slicez"
)

const builtinPkgPath = "ezpkg.io/codez/builtin"

type Packages struct {
	Fset *token.FileSet

	pkgs []*Package

	mapPkgs  map[string]*Package
	allPkgs  []*Package
	stdPkgs  []*Package
	pkgByPos []*Package
	builtin  map[string]types.Type
}

type Package struct {
	*packages.Package

	start, end token.Pos
}

func newPackage(pkg *packages.Package) *Package {
	p := &Package{
		Package: pkg,
	}
	for _, file := range pkg.Syntax {
		if p.start == 0 || file.Pos() < p.start {
			p.start = file.Pos()
		}
		if p.end == 0 || file.End() > p.end {
			p.end = file.End()
		}
	}
	return p
}

func (p *Package) GetObject(name string) types.Object {
	return p.Types.Scope().Lookup(name)
}
func (p *Package) GetType(obj types.Object) types.Type {
	return obj.Type()
}
func (p *Package) Positions() (start, end token.Pos) {
	return p.start, p.end
}
func (p *Package) HasPos(pos token.Pos) bool {
	if pos < p.start || pos > p.end {
		return false
	}
	for _, file := range p.Syntax {
		if file.Pos() <= pos && pos <= file.End() {
			return true
		}
	}
	return false
}
func (p *Package) quickHasPos(pos token.Pos) bool {
	return p.start <= pos && pos <= p.end
}

// Packages returns the loaded packages from input patterns.
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

func (p *Packages) GetPackageByPath(path string) *Package {
	return p.mapPkgs[path]
}

func (p *Packages) GetObject(pkgPath, objName string) types.Object {
	if pkgPath == "" {
		return types.Universe.Lookup(objName)
	}
	pkg := p.GetPackageByPath(pkgPath)
	if pkg == nil {
		return nil
	}
	return pkg.GetObject(objName)
}

func (p *Packages) MustGetObject(pkgPath, objName string) types.Object {
	obj := p.GetObject(pkgPath, objName)
	if obj == nil {
		panic(fmt.Sprintf("object %v.%v not found", pkgPath, objName))
	}
	return obj
}

func (p *Packages) GetType(pkgPath, objName string) types.Type {
	obj := p.GetObject(pkgPath, objName)
	if obj == nil {
		return nil
	}
	return obj.Type()
}

func (p *Packages) MustGetType(pkgPath, objName string) types.Type {
	typ := p.GetType(pkgPath, objName)
	if typ == nil {
		panic(fmt.Sprintf("type %v.%v not found", pkgPath, objName))
	}
	return typ
}

func (p *Packages) GetBuiltInType(typName string) types.Type {
	obj := types.Universe.Lookup(typName)
	if obj == nil {
		return nil
	}
	return obj.Type()
}

func (p *Packages) MustGetBuiltInType(typName string) types.Type {
	typ := p.GetBuiltInType(typName)
	if typ == nil {
		panic(fmt.Sprintf("type %q not found", typName))
	}
	return typ
}

func (p *Packages) GetPackageByPos(pos token.Pos) *Package {
	pkg := quickSearchPkgByPos(p.pkgByPos, pos)
	if pkg != nil && pkg.HasPos(pos) {
		return pkg
	}
	// slow path
	for _, pkg := range p.pkgs {
		if pkg.HasPos(pos) {
			return pkg
		}
	}
	return nil
}

func quickSearchPkgByPos(pkgs []*Package, pos token.Pos) *Package {
	switch {
	case len(pkgs) == 0:
		return nil
	case len(pkgs) == 1 && pkgs[0].quickHasPos(pos):
		return pkgs[0]
	case len(pkgs) == 1:
		return nil
	}
	mid := len(pkgs) / 2
	if pkgs[mid].quickHasPos(pos) {
		return pkgs[mid]
	}
	if pkgs[mid].start > pos {
		return quickSearchPkgByPos(pkgs[:mid], pos)
	} else {
		return quickSearchPkgByPos(pkgs[mid:], pos)
	}
}

func (p *Packages) TypeOf(expr ast.Expr) types.Type {
	pkg := p.GetPackageByPos(expr.Pos())
	if pkg == nil {
		return nil
	}
	return pkg.TypesInfo.TypeOf(expr)
}

func (p *Packages) ObjectOf(ident *ast.Ident) types.Object {
	for _, pkg := range p.pkgs {
		if obj := pkg.TypesInfo.ObjectOf(ident); obj != nil {
			return obj
		}
	}
	return nil
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
	sortPkgsByPos := func(pkgs []*Package) []*Package {
		slices.SortFunc(pkgs, func(a, b *Package) int {
			return cmp.Compare(a.start, b.start)
		})
		return pkgs
	}

	if len(pkgs) == 0 {
		return nil
	}
	var allPkgs []*Package
	p := &Packages{pkgs: pkgs, Fset: pkgs[0].Fset}
	p.mapPkgs = map[string]*Package{}
	for _, pkg := range p.pkgs {
		for path, impPkg := range pkg.Imports {
			if p.mapPkgs[path] == nil {
				pkg0 := newPackage(impPkg)
				p.mapPkgs[path] = pkg0
				allPkgs = append(allPkgs, pkg0)
			}
		}
	}

	// filter std packages then sort
	var goOrgPkgs, otherPkgs []*Package
	for _, pkg := range allPkgs {
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
	p.pkgByPos = sortPkgsByPos(allPkgs)
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
