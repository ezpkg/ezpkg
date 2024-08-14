package ggen

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/tools/go/packages"

	"ezpkg.io/errorz"
)

type Positioner interface {
	Pos() token.Pos
}

type GeneratingPackage struct {
	*packages.Package

	directives []Directive
	plugin     *pluginStruct
	engine     *engine
	printer    *printer
}

func (g *GeneratingPackage) GetDir() string {
	dir := getPackageDir(g.Package)
	if dir == "" {
		panic(fmt.Sprintf("can not get directory of package %q", g.PkgPath))
	}
	return dir
}

func (g *GeneratingPackage) GetPrinter() Printer {
	if g.printer == nil {
		fileName := generateFileName(g.engine, g.plugin)
		filePath := filepath.Join(getPackageDir(g.Package), fileName)
		g.printer = newPrinter(g.engine, g.plugin, g.Package.Types, "", filePath)
	}
	return g.printer
}

func (g *GeneratingPackage) GetDirectives() []Directive {
	return cloneDirectives(g.directives)
}

func (g *GeneratingPackage) GetObjects() []types.Object {
	return g.engine.GetObjectsByPackage(g.Package)
}

type Engine interface {

	// Plugin should use the embedded logger to log messages.
	Logger

	// GenerateEachPackage loops through the list of GeneratingPackages and call the given function.
	GenerateEachPackage(func(Engine, *packages.Package, Printer) error) error

	// GeneratingPackages returns a list of packages available for generating.
	GeneratingPackages() []*GeneratingPackage

	// GeneratePackage generates file at given package path with the given file name. The file name must not include any slash character (/). If fileName is empty, use default file name.
	GeneratePackage(pkg *packages.Package, fileName string) (Printer, error)

	// GenerateFile generates file at given path. It should be an absolute path, can include slash character (/). If the path ends with /, use default file name.
	GenerateFile(pkgName, filePath string) (Printer, error)

	GetComment(Positioner) Comment
	GetDirectives(Positioner) Directives
	GetDirectivesByPackage(*packages.Package) Directives
	GetIdent(Positioner) *ast.Ident
	GetObject(Positioner) types.Object
	GetObjectByName(pkgPath, name string) types.Object
	GetBuiltinType(name string) types.Type
	GetObjectsByPackage(*packages.Package) []types.Object
	GetObjectsByScope(*types.Scope) []types.Object
	GetPackage(Positioner) *packages.Package
	GetPackageByPath(string) *packages.Package

	LogDebugNode(node ast.Node) error
}

var _ Engine = &wrapEngine{}

type engine struct {
	logger Logger

	plugins        []*pluginStruct
	enabledPlugins []*pluginStruct
	pluginsMap     map[string]*pluginStruct

	xcfg    Config
	xinfo   *extendedInfo
	pkgcfg  packages.Config
	pkgMap  map[string]*packages.Package
	dir2pkg map[string]*packages.Package
	srcMap  map[string][]byte
	bufPool *sync.Pool

	builtinTypes           map[string]types.Type
	cleanedFileNames       map[string]bool
	mapPkgDirectives       map[string][]Directive
	collectedPackages      []filteringPackage
	includedPatterns       []string
	includedPackages       map[string][]bool
	sortedIncludedPackages []includedPackage
	generatedFiles         []string
}

type wrapEngine struct {
	zLogger
	*engine

	plugin *pluginStruct
	pkgs   []*GeneratingPackage
}

func newEngine(logger Logger) *engine {
	return &engine{
		logger:     logger,
		pkgMap:     make(map[string]*packages.Package),
		dir2pkg:    make(map[string]*packages.Package),
		pluginsMap: make(map[string]*pluginStruct),
		bufPool:    &sync.Pool{},
	}
}

func (ng *engine) GetComment(p Positioner) Comment {
	cmt := ng.xinfo.GetComment(ng.GetIdent(p))
	return cmt
}

func (ng *engine) CommentByIdent(ident *ast.Ident) Comment {
	cmt := ng.xinfo.GetComment(ident)
	return cmt
}

func (ng *engine) CommentByObject(obj types.Object) Comment {
	ident := ng.GetIdentByPos(obj.Pos())
	return ng.CommentByIdent(ident)
}

func (ng *engine) GetDirectives(p Positioner) Directives {
	return ng.GetComment(p).Directives
}

func (ng *engine) GetIdent(p Positioner) *ast.Ident {
	return ng.GetIdentByPos(p.Pos())
}

func (ng *engine) GetIdentByObject(obj types.Object) *ast.Ident {
	return ng.GetIdentByPos(obj.Pos())
}

func (ng *engine) GetIdentByPos(pos token.Pos) *ast.Ident {
	return ng.xinfo.Positions[pos]
}

func (ng *engine) GetObject(p Positioner) types.Object {
	return ng.GetObjectByIdent(ng.GetIdent(p))
}

func (ng *engine) GetObjectByIdent(ident *ast.Ident) types.Object {
	return ng.xinfo.GetObject(ident)
}

func (ng *engine) GetObjectByName(pkgPath, name string) types.Object {
	pkg := ng.GetPackageByPath(pkgPath)
	if pkg == nil {
		return nil
	}
	return pkg.Types.Scope().Lookup(name)
}

func (ng *engine) GetBuiltinType(name string) types.Type {
	return ng.builtinTypes[name]
}

func (ng *engine) GetPackage(p Positioner) *packages.Package {
	return ng.GetPackageByIdent(ng.GetIdent(p))
}

func (ng *engine) GetPackageByIdent(ident *ast.Ident) *packages.Package {
	decl := ng.xinfo.Declarations[ident]
	if decl == nil {
		return nil
	}
	return decl.Pkg
}

func (ng *engine) GetPackageByPath(pkgPath string) *packages.Package {
	return ng.pkgMap[pkgPath]
}

func (ng *engine) GetObjectsByPackage(pkg *packages.Package) []types.Object {
	return ng.GetObjectsByScope(pkg.Types.Scope())
}

func (ng *engine) GetObjectsByScope(s *types.Scope) []types.Object {
	names := s.Names()
	objs := make([]types.Object, len(names))
	for i, name := range names {
		objs[i] = s.Lookup(name)
	}
	return objs
}

func (ng *wrapEngine) Logger() Logger {
	return ng.logger
}

func (ng *wrapEngine) GenerateEachPackage(
	fn func(Engine, *packages.Package, Printer) error,
) error {
	for _, pkg := range ng.generatingPackages() {
		prt := pkg.GetPrinter()
		if err := fn(ng, pkg.Package, prt); err != nil {
			return errorz.Wrapf(err, "generating package %v: %v", pkg.PkgPath, err)
		}
		if len(prt.Bytes()) == 0 {
			continue
		}
		if err := prt.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (ng *wrapEngine) GeneratingPackages() []*GeneratingPackage {
	if ng.pkgs == nil {
		ng.pkgs = ng.generatingPackages()
	}
	return ng.pkgs
}

func (ng *wrapEngine) generatingPackages() []*GeneratingPackage {
	index := ng.plugin.index
	pkgs := make([]*GeneratingPackage, 0, len(ng.engine.sortedIncludedPackages))
	for _, p := range ng.engine.sortedIncludedPackages {
		if p.Included != nil && p.Included[index] {
			pkg := ng.engine.pkgMap[p.PkgPath]
			if pkg == nil {
				continue
			}
			gpkg := ng.generatingPackage(pkg)
			pkgs = append(pkgs, gpkg)
		}
	}
	return pkgs
}

func (ng *wrapEngine) generatingPackage(pkg *packages.Package) *GeneratingPackage {
	directives := ng.GetDirectivesByPackage(pkg)
	gpkg := &GeneratingPackage{
		Package:    pkg,
		directives: directives,
		plugin:     ng.plugin,
		engine:     ng.engine,
	}
	return gpkg
}

func generateFileName(ng *engine, plugin *pluginStruct) string {
	input := GenerateFileNameInput{PluginName: plugin.name}
	return ng.genFilename(input)
}

func (ng *wrapEngine) GeneratePackage(pkg *packages.Package, fileName string) (Printer, error) {
	if strings.Contains(fileName, "/") {
		return nil, errorz.Wrapf(nil, "invalid filename: file must not contain / (filename=%v)", fileName)
	}
	if fileName == "" {
		fileName = generateFileName(ng.engine, ng.plugin)
	}
	filePath := filepath.Join(getPackageDir(pkg), fileName)
	prt := newPrinter(ng.engine, ng.plugin, pkg.Types, "", filePath)
	return prt, nil
}

func (ng *wrapEngine) GenerateFile(pkgName string, filePath string) (Printer, error) {
	var pkg *types.Package
	if filePath == "" {
		return nil, errorz.Wrapf(nil, "empty file path")
	}
	dir := filepath.Dir(filePath)
	if pkg0 := ng.dir2pkg[dir]; pkg0 != nil {
		pkg = pkg0.Types
		pkgName = pkg0.Name
	}
	if pkgName == "" {
		return nil, errorz.Wrapf(nil, "empty package name")
	}

	if strings.HasSuffix(filePath, "/") {
		fileName := generateFileName(ng.engine, ng.plugin)
		filePath = filepath.Join(filePath, fileName)
	}
	{
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return nil, errorz.Wrapf(err, "create directory %q: %v", dir, err)
		}
		file, err := os.Open(dir)
		if err != nil {
			return nil, errorz.Wrapf(err, "can not read dir %q: %v", dir, err)
		}
		names, err := file.Readdirnames(-1)
		if err != nil {
			return nil, errorz.Wrapf(err, "can not read dir %q: %v", dir, err)
		}
		found := false
		for _, name := range names {
			if strings.HasSuffix(name, ".go") {
				found = true
				break
			}
		}
		if !found {
			// create an empty doc.go for working around "can not find module
			// providing package ..." error
			docFile := filepath.Join(dir, "doc.go")
			err = os.WriteFile(docFile, []byte("package "+pkgName), 0644)
			if err != nil {
				return nil, errorz.Wrapf(err, "can not write file %v: %v", docFile, err)
			}
		}
	}
	pr := newPrinter(ng.engine, ng.plugin, pkg, pkgName, filePath)
	return pr, nil
}

func (ng *wrapEngine) GetDirectivesByPackage(pkg *packages.Package) Directives {
	directives, ok := ng.engine.mapPkgDirectives[pkg.PkgPath]
	if !ok {
		for _, file := range pkg.GoFiles {
			body, err := os.ReadFile(file)
			if err != nil {
				if os.IsNotExist(err) {
					ng.Errorw("ignore not found file", "file", file)
					continue
				}
				panic(err)
			}

			// only parse package level directives
			errs := parseDirectivesFromBody(body, &directives, nil)
			for _, err = range errs {
				ng.Errorw("invalid directive from file", "error", err, "file", file)
			}
		}
		ng.engine.mapPkgDirectives[pkg.PkgPath] = directives
	}
	return cloneDirectives(directives)
}

func (ng *wrapEngine) LogDebugNode(node ast.Node) error {
	return ast.Print(ng.engine.xinfo.Fset, node)
}

func cloneDirectives(directives []Directive) []Directive {
	if len(directives) == 0 {
		return nil
	}
	result := make([]Directive, len(directives))
	copy(result, directives)
	return result
}
