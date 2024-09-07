package codez

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"iter"
	"slices"
	"strings"

	"golang.org/x/tools/go/packages"

	"ezpkg.io/errorz"
	"ezpkg.io/iterz"
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

// wrapPackage wraps packages.Package to PackageX, and store it in Packages.mapPkgs.
func (px *Packages) wrapPackage(pkg *packages.Package) *PackageX {
	if p := px.mapPkgs[pkg.PkgPath]; p != nil {
		return p
	}
	p := &PackageX{Package: pkg, px: px}
	px.mapPkgs[pkg.PkgPath] = p
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

func (p *PackageX) GetFileByName(name string) *FileX {
	for _, file := range p.Syntax {
		if file.Name.Name == name {
			tokFile := p.Fset.File(file.Pos())
			return p.px.newFileX(file, tokFile, p)
		}
	}
	return nil
}

func (p *PackageX) MustGetFileByName(name string) *FileX {
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

func newPackages(pkgs []*packages.Package) (*Packages, error) {
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
		return nil, errorz.New("no packages loaded")
	}

	px := &Packages{Fset: pkgs[0].Fset}
	px.mapPkgs = map[string]*PackageX{}
	for _, pkg := range pkgs {
		slicez.AppendTo(&px.origPkgs, px.wrapPackage(pkg))
		for _, impPkg := range pkg.Imports {
			px.wrapPackage(impPkg)
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
	return px, nil
}

// Packages returns the loaded packages from input patterns.
func (px *Packages) InputPackages() []*PackageX {
	return px.origPkgs
}

// AllPackages returns all packages, including std packages, golang.org/x packages, and other packages. It supports filtering by pattern. Examples:
//
//	AllPackages()                         👉 return all packages
//	AllPackages("ezpkg.io/...")           👉 return packages that start with "ezpkg.io"
//	AllPackages("ezpkg.io/codez", "fmt")  👉 return listed packages
func (px *Packages) AllPackages(pattern ...string) []*PackageX {
	return filterPackages(px.allPkgs, pattern...)
}
func (px *Packages) AllErrors() (errs Errors) {
	for _, pkg := range px.origPkgs {
		errs = append(errs, pkg.Errors...)
	}
	return errs
}

// FirstErrors collect the first error message for each error package.
func (px *Packages) FirstErrors() (errs Errors) {
	if px == nil {
		return nil
	}
	for _, pkg := range px.origPkgs {
		if len(pkg.Errors) > 0 {
			errs = append(errs, pkg.Errors[0])
		}
	}
	return errs
}

func (px *Packages) StdPackages() []*PackageX {
	return px.stdPkgs
}
func (px *Packages) NonStdPackages() []*PackageX {
	return px.allPkgs[len(px.stdPkgs):]
}

func (px *Packages) GetPackageByPath(path string) *PackageX {
	return px.mapPkgs[path]
}

func (px *Packages) MustGetPackageByPath(path string) *PackageX {
	pkg := px.GetPackageByPath(path)
	if pkg == nil {
		panic(fmt.Sprintf("package %q not found", path))
	}
	return pkg
}

func (px *Packages) GetObject(pkgPath, objName string) types.Object {
	if pkgPath == "" {
		return types.Universe.Lookup(objName)
	}
	pkg := px.GetPackageByPath(pkgPath)
	if pkg == nil {
		return nil
	}
	return pkg.GetObject(objName)
}

func (px *Packages) MustGetObject(pkgPath, objName string) types.Object {
	obj := px.GetObject(pkgPath, objName)
	if obj == nil {
		panic(fmt.Sprintf("object %v.%v not found", pkgPath, objName))
	}
	return obj
}

func (px *Packages) GetType(pkgPath, objName string) types.Type {
	obj := px.GetObject(pkgPath, objName)
	if obj == nil {
		return nil
	}
	return obj.Type()
}

func (px *Packages) MustGetType(pkgPath, objName string) types.Type {
	typ := px.GetType(pkgPath, objName)
	if typ == nil {
		panic(fmt.Sprintf("type %v.%v not found", pkgPath, objName))
	}
	return typ
}

func (px *Packages) GetBuiltInType(typName string) types.Type {
	obj := types.Universe.Lookup(typName)
	if obj == nil {
		return nil
	}
	return obj.Type()
}

func (px *Packages) MustGetBuiltInType(typName string) types.Type {
	typ := px.GetBuiltInType(typName)
	if typ == nil {
		panic(fmt.Sprintf("type %q not found", typName))
	}
	return typ
}

func (px *Packages) GetPackageByPos(pos token.Pos) *PackageX {
	// TODO: optimize this
	for _, pkg := range px.allPkgs {
		if pkg.HasPos(pos) {
			return pkg
		}
	}
	return nil
}

func (px *Packages) MustGetPackageByPos(pos token.Pos) *PackageX {
	pkg := px.GetPackageByPos(pos)
	if pkg == nil {
		panic(fmt.Sprintf("package not found at %v", pos))
	}
	return pkg
}

func (px *Packages) GetFileByPos(pos token.Pos) *FileX {
	pkg := px.GetPackageByPos(pos)
	if pkg == nil {
		return nil
	}
	return pkg.GetFileByPos(pos)
}

func (px *Packages) MustGetFileByPos(pos token.Pos) *FileX {
	file := px.GetFileByPos(pos)
	if file == nil {
		panic(fmt.Sprintf("file not found at %v", pos))
	}
	return file
}

func (px *Packages) TypeOf(expr ast.Expr) types.Type {
	if t, ok := px.Types[expr]; ok {
		return t.Type
	}
	if id, _ := expr.(*ast.Ident); id != nil {
		if obj := px.ObjectOf(id); obj != nil {
			return obj.Type()
		}
	}
	return nil
}

func (px *Packages) ObjectOf(ident *ast.Ident) types.Object {
	if obj := px.Defs[ident]; obj != nil {
		return obj
	}
	return px.Uses[ident]
}

func (px *Packages) PkgNameOf(imp *ast.ImportSpec) *types.PkgName {
	var obj types.Object
	if imp.Name != nil {
		obj = px.Defs[imp.Name]
	} else {
		obj = px.Implicits[imp]
	}
	pkgname, _ := obj.(*types.PkgName)
	return pkgname
}

func (px *Packages) IterNodesByPos(pos token.Pos) iter.Seq2[NodePath, *NodeX] {
	file := px.GetFileByPos(pos)
	if file == nil {
		return iterz.Nil2[NodePath, *NodeX]()
	}
	return file.IterNodesByPos(pos)
}

func (px *Packages) MustNodesByPos(pos token.Pos) iter.Seq2[NodePath, *NodeX] {
	file := px.MustGetFileByPos(pos)
	return file.MustNodesByPos(pos)
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
