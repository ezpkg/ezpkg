package codez

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"slices"
	"strings"

	"golang.org/x/tools/go/packages"

	"ezpkg.io/mapz"
	"ezpkg.io/slicez"
)

type Packages struct {
	Fset *token.FileSet

	origPkgs []*PackageX          // original packages from input patterns
	mapPkgs  map[string]*PackageX // map of all packages by path
	allPkgs  []*PackageX          // all packages, including std packages
	stdPkgs  []*PackageX          // std packages

	// --- collect types.Info from all packages ---

	Types      map[ast.Expr]types.TypeAndValue
	Instances  map[*ast.Ident]types.Instance
	Defs       map[*ast.Ident]types.Object
	Uses       map[*ast.Ident]types.Object
	Implicits  map[ast.Node]types.Object
	Selections map[*ast.SelectorExpr]*types.Selection
	Scopes     map[ast.Node]*types.Scope
}

type PackageX struct {
	*packages.Package
	px *Packages
}

func (px *Packages) newPackage(pkg *packages.Package) *PackageX {
	p := &PackageX{
		Package: pkg,
		px:      px,
	}
	return p
}

func (p *PackageX) Unwrap() *packages.Package { return p.Package }

func (p *PackageX) GetObject(name string) types.Object {
	return p.Types.Scope().Lookup(name)
}

func (p *PackageX) MustGetObject(name string) types.Object {
	obj := p.GetObject(name)
	if obj == nil {
		panic(fmt.Sprintf("object %q not found", name))
	}
	return obj
}

func (p *PackageX) GetFileByPos(pos token.Pos) *FileX {
	for _, file := range p.Syntax {
		if file.FileStart <= pos && pos <= file.FileEnd {
			tokFile := p.Fset.File(file.Pos())
			return p.px.newFileX(file, tokFile, p)
		}
	}
	return nil
}

func (p *PackageX) MustGetFileByPos(pos token.Pos) *FileX {
	file := p.GetFileByPos(pos)
	if file == nil {
		panic(fmt.Sprintf("file not found at %v", pos))
	}
	return file
}

func (p *PackageX) GetFileByName(name string) *ast.File {
	for _, file := range p.Syntax {
		if file.Name.Name == name {
			return file
		}
	}
	return nil
}

func (p *PackageX) MustGetFileByName(name string) *ast.File {
	file := p.GetFileByName(name)
	if file == nil {
		panic(fmt.Sprintf("file %q not found", name))
	}
	return file
}

func (p *PackageX) HasPos(pos token.Pos) bool {
	for _, file := range p.Syntax {
		if file.FileStart <= pos && pos <= file.FileEnd {
			return true
		}
	}
	return false
}

func newPackages(pkgs []*PackageX) *Packages {
	isStd := func(path string) bool {
		return !strings.Contains(path, ".")
	}
	isGolangOrg := func(path string) bool {
		return strings.HasPrefix(path, "golang.org/")
	}
	sortPkgs := func(pkgs []*PackageX) {
		slices.SortFunc(pkgs, func(a, b *PackageX) int {
			return strings.Compare(a.PkgPath, b.PkgPath)
		})
	}

	if len(pkgs) == 0 {
		return nil
	}

	px := &Packages{origPkgs: pkgs, Fset: pkgs[0].Fset}
	px.mapPkgs = map[string]*PackageX{}
	for _, pkg := range px.origPkgs {
		px.mapPkgs[pkg.PkgPath] = pkg
		for path, impPkg := range pkg.Imports {
			if px.mapPkgs[path] == nil {
				pkg0 := px.newPackage(impPkg)
				px.mapPkgs[path] = pkg0
			}
		}
	}
	_, listPkgs := mapz.SortedKeysAndValues(px.mapPkgs)

	// filter std packages then sort
	var goOrgPkgs, otherPkgs []*PackageX
	for _, pkg := range listPkgs {
		switch {
		case isStd(pkg.PkgPath):
			px.stdPkgs = append(px.stdPkgs, pkg)
		case isGolangOrg(pkg.PkgPath):
			goOrgPkgs = append(goOrgPkgs, pkg)
		default:
			otherPkgs = append(otherPkgs, pkg)
		}
	}
	sortPkgs(px.origPkgs)
	sortPkgs(px.stdPkgs)
	sortPkgs(goOrgPkgs)
	sortPkgs(otherPkgs)
	px.allPkgs = slicez.Concat(px.stdPkgs, goOrgPkgs, otherPkgs)

	// collect types.Info from all packages
	px.Types = map[ast.Expr]types.TypeAndValue{}
	px.Instances = map[*ast.Ident]types.Instance{}
	px.Defs = map[*ast.Ident]types.Object{}
	px.Uses = map[*ast.Ident]types.Object{}
	px.Implicits = map[ast.Node]types.Object{}
	px.Selections = map[*ast.SelectorExpr]*types.Selection{}
	px.Scopes = map[ast.Node]*types.Scope{}
	for _, pkg := range listPkgs {
		mapz.Append(px.Types, pkg.TypesInfo.Types)
		mapz.Append(px.Instances, pkg.TypesInfo.Instances)
		mapz.Append(px.Defs, pkg.TypesInfo.Defs)
		mapz.Append(px.Uses, pkg.TypesInfo.Uses)
		mapz.Append(px.Implicits, pkg.TypesInfo.Implicits)
		mapz.Append(px.Selections, pkg.TypesInfo.Selections)
		mapz.Append(px.Scopes, pkg.TypesInfo.Scopes)
	}
	return px
}

// Packages returns the loaded packages from input patterns.
func (p *Packages) InputPackages() []*PackageX {
	return p.origPkgs
}

// AllPackages returns all packages, including std packages, golang.org/x packages, and other packages. It supports filtering by pattern. Examples:
//
//	AllPackages()                         👉 return all packages
//	AllPackages("ezpkg.io/...")           👉 return packages that start with "ezpkg.io"
//	AllPackages("ezpkg.io/codez", "fmt")  👉 return listed packages
func (p *Packages) AllPackages(pattern ...string) []*PackageX {
	return filterPackages(p.allPkgs, pattern...)
}
func (p *Packages) AllErrors() (errs []packages.Error) {
	for _, pkg := range p.origPkgs {
		errs = append(errs, pkg.Errors...)
	}
	return errs
}
func (p *Packages) StdPackages() []*PackageX {
	return p.stdPkgs
}
func (p *Packages) NonStdPackages() []*PackageX {
	return p.allPkgs[len(p.stdPkgs):]
}

func (p *Packages) GetPackageByPath(path string) *PackageX {
	return p.mapPkgs[path]
}

func (p *Packages) MustGetPackageByPath(path string) *PackageX {
	pkg := p.GetPackageByPath(path)
	if pkg == nil {
		panic(fmt.Sprintf("package %q not found", path))
	}
	return pkg
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

func (p *Packages) GetPackageByPos(pos token.Pos) *PackageX {
	// TODO: optimize this
	for _, pkg := range p.allPkgs {
		if pkg.HasPos(pos) {
			return pkg
		}
	}
	return nil
}

func (p *Packages) MustGetPackageByPos(pos token.Pos) *PackageX {
	pkg := p.GetPackageByPos(pos)
	if pkg == nil {
		panic(fmt.Sprintf("package not found at %v", pos))
	}
	return pkg
}

func (p *Packages) GetFileByPos(pos token.Pos) *FileX {
	pkg := p.GetPackageByPos(pos)
	if pkg == nil {
		return nil
	}
	return pkg.GetFileByPos(pos)
}

func (p *Packages) MustGetFileByPos(pos token.Pos) *FileX {
	file := p.GetFileByPos(pos)
	if file == nil {
		panic(fmt.Sprintf("file not found at %v", pos))
	}
	return file
}

func (p *Packages) TypeOf(expr ast.Expr) types.Type {
	if t, ok := p.Types[expr]; ok {
		return t.Type
	}
	if id, _ := expr.(*ast.Ident); id != nil {
		if obj := p.ObjectOf(id); obj != nil {
			return obj.Type()
		}
	}
	return nil
}

func (p *Packages) ObjectOf(ident *ast.Ident) types.Object {
	if obj := p.Defs[ident]; obj != nil {
		return obj
	}
	return p.Uses[ident]
}

func (p *Packages) PkgNameOf(imp *ast.ImportSpec) *types.PkgName {
	var obj types.Object
	if imp.Name != nil {
		obj = p.Defs[imp.Name]
	} else {
		obj = p.Implicits[imp]
	}
	pkgname, _ := obj.(*types.PkgName)
	return pkgname
}

func filterPackages(pkgs []*PackageX, pattern ...string) []*PackageX {
	if len(pattern) == 0 {
		return pkgs
	}
	var out []*PackageX
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
