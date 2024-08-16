package main

import (
	"fmt"
	"go/types"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
	"golang.org/x/tools/go/packages"

	"ezpkg.io/errorz"
	"ezpkg.io/genz"
	"ezpkg.io/mapz"
	"ezpkg.io/slicez"
	"ezpkg.io/typez"
)

const (
	pathCodez = "ezpkg.io/codez"
	pathGoAst = "go/ast"
)

func cmdCodezMatchers(cx *cli.Context) error {
	cfg := initConfig(cx, newCodezMatcher())
	return genz.Start(cx.Context, cfg, "ezpkg.io/codez")
}

type CodezMatcher struct {
}

func newCodezMatcher() genz.Plugin {
	return &CodezMatcher{}
}

func (c CodezMatcher) Name() string {
	return "codez-matchers"
}

func (c CodezMatcher) Filter(ng genz.FilterEngine) error {
	ng.IncludePackage(pathCodez)
	return nil
}

func (c CodezMatcher) Generate(ng genz.Engine) error {
	pkgCodez := ng.GetPackageByPath(pathCodez)
	pkgDir := filepath.Dir(pkgCodez.CompiledGoFiles[0])

	pkgGoAst := ng.GetPackageByPath(pathGoAst)
	allObjs := getObjs(pkgGoAst)
	allObjs = slicez.FilterFunc(allObjs, func(obj types.Object) bool {
		return !typez.In(obj.Name(), "BadDecl", "BadExpr", "BadStmt")
	})

	_, astNodeI := getIface(pkgGoAst, "Node")
	_, astExprI := getIface(pkgGoAst, "Expr")
	_, astStmtI := getIface(pkgGoAst, "Stmt")
	_, astDeclI := getIface(pkgGoAst, "Decl")

	matchers := map[string]types.Type{}
	addMatcher := func(obj types.Object) {
		if matchers[obj.Name()] == nil {
			matchers[obj.Name()] = obj.Type()
		}
	}

	type ParsedField struct {
		isNode bool
		token  *types.Named
		astTyp *types.Named
		slice  *types.Slice
		basic  *types.Basic
	}
	parseField := func(field types.Type) *ParsedField {
		if named := asNamed(field); named != nil {
			if typez.In(named.Obj().Name(), "Object") {
				return nil // ignore deprecated types
			}
		}
		astTyp := asAstType(field)
		return &ParsedField{
			isNode: astTyp != nil && implements(astTyp, astNodeI),
			token:  asTokenType(field),
			astTyp: astTyp,
			slice:  asSlice(field),
			basic:  asBasic(field),
		}
	}
	gen := func(class string, iface *types.Interface) {
		Class := title(class)
		file := fmt.Sprintf("matchers.%s.go", class)
		p := errorz.Must(ng.GenerateFile("codez", pkgDir+"/"+file))
		p.Import("ast", "go/ast")
		defer func() { errorz.MustZ(p.Close()) }()
		pr := p.Printf

		for _, x := range allObjs.Implements(iface).Structs() {
			addMatcher(x)
			zName := fmt.Sprintf("z%sMatcher", x.Name())

			pr("// %s\n", x.Name())
			pr("type %s struct {\n", zName)
			pr("\t_ *%s\n\n", p.TypeString(x.Type()))
			st := mustStruct(x.Type())
			for i := 0; i < st.NumFields(); i++ {
				field := st.Field(i)
				f := parseField(field.Type())
				if f == nil {
					continue
				}

				switch {
				case f.token != nil:
					pr("\t%s %s\n", field.Name(), p.TypeString(field.Type()))

				case f.astTyp != nil:
					pr("\t%s %sMatcher\n", field.Name(), f.astTyp.Obj().Name())

				case f.slice != nil:
					typ := asAstType(f.slice.Elem())
					if typ == nil {
						panic(fmt.Sprintf("unsupported slice type %v", field.Type()))
					}
					pr("\t%s %sListMatcher[ast.%s]\n", field.Name(), typ.Obj().Name(), typ.Obj().Name())

				case f.basic != nil:
					pr("\t%s %sMatcher\n", field.Name(), title(field.Type().String()))

				default:
					pr("\t%s %s ❌\n", field.Name(), p.TypeString(field.Type()))
				}
			}
			pr("}\n\n")

			pr("func (m %s) Match%s(node ast.%s) (ok bool, err error) {\n", zName, Class, Class)
			pr("\treturn m.Match(node)\n")
			pr("}\n")

			pr("func (m %s) Match(node ast.Node) (ok bool, err error) {\n", zName)
			pr("\tx, ok := node.(*ast.%s)\n", x.Name())
			pr("\tif !ok {\n")
			pr("\t\treturn false, nil\n")
			pr("\t}\n")

			for i := 0; i < st.NumFields(); i++ {
				field := st.Field(i)
				f := parseField(field.Type())
				if f == nil {
					continue
				}

				switch {
				case f.token != nil:
					continue
				case f.basic != nil:
					pr("\tok, err = matchValue(ok, err, m.%s, x.%s)\n", field.Name(), field.Name())
				case f.isNode:
					pr("\tok, err = match(ok, err, m.%s, x.%s)\n", field.Name(), field.Name())
				case f.slice != nil:
					pr("\tok, err = matchList(ok, err, m.%s, x.%s)\n", field.Name(), field.Name())
				default:
					pr("\tok, err = matchValue(ok, err, m.%s, x.%s)\n", field.Name(), field.Name())
				}
			}
			pr("\treturn ok, err\n")
			pr("}\n\n")
		}
	}
	gen("expr", astExprI)
	gen("stmt", astStmtI)
	gen("decl", astDeclI)

	{ // 👉 interfaces
		p := errorz.Must(ng.GenerateFile("codez", pkgDir+"/matchers.iface.go"))
		p.Import("ast", "go/ast")
		defer func() { errorz.MustZ(p.Close()) }()
		pr := p.Printf

		for _, name := range mapz.SortedKeys(matchers) {
			matcher := matchers[name]

			pr("// %v\n", name)
			pr("type %sMatcher interface {\n", name)
			pr("\tMatch(node ast.Node) (bool, error)\n")
			if implements(matcher, astExprI) {
				pr("\tMatchExpr(expr ast.Expr) (bool, error)\n")
			}
			if implements(matcher, astStmtI) {
				pr("\tMatchStmt(stmt ast.Stmt) (bool, error)\n")
			}
			if implements(matcher, astDeclI) {
				pr("\tMatchDecl(decl ast.Decl) (bool, error)\n")
			}
			pr("}\n\n")
		}
	}
	return nil
}

type Objects []types.Object

func (objs Objects) Implements(iface *types.Interface) (out Objects) {
	for _, obj := range objs {
		if implements(obj.Type(), iface) {
			out = append(out, obj)
		}
	}
	return out
}

func (objs Objects) Structs() (out Objects) {
	for _, obj := range objs {
		if _, ok := obj.Type().Underlying().(*types.Struct); ok {
			out = append(out, obj)
		}
	}
	return out
}

func implements(typ types.Type, iface *types.Interface) bool {
	ptr := types.NewPointer(typ)
	return types.Implements(typ, iface) || types.Implements(ptr, iface)
}

func getObj(pkg *types.Package, objName string) types.Object {
	return pkg.Scope().Lookup(objName)
}

func getObjs(pkg *packages.Package) Objects {
	scope := pkg.Types.Scope()
	return slicez.MapFunc(scope.Names(), func(name string) types.Object {
		return scope.Lookup(name)
	})
}

func getIface(pkg *packages.Package, iface string) (types.Object, *types.Interface) {
	obj := pkg.Types.Scope().Lookup(iface)
	typ := obj.Type().(*types.Named).Underlying().(*types.Interface)
	return obj, typ
}

func asNamed(typ types.Type) *types.Named {
	typ = skipPtr(typ)
	named, _ := typ.(*types.Named)
	return named
}

func asStruct(typ types.Type) *types.Struct {
	st, _ := asNamed(typ).Underlying().(*types.Struct)
	return st
}

func skipPtr(typ types.Type) types.Type {
	ptr, _ := typ.(*types.Pointer)
	if ptr != nil {
		return ptr.Elem()
	}
	return typ
}

func skipNamed(typ types.Type) types.Type {
	named, _ := typ.(*types.Named)
	if named != nil {
		return named.Underlying()
	}
	return typ
}

func mustStruct(typ0 types.Type) *types.Struct {
	typ := skipPtr(typ0)
	typ = skipNamed(typ)
	st, _ := typ.(*types.Struct)
	if st == nil {
		panic(fmt.Sprintf("%v is not a struct", typ0))
	}
	return st
}

func asTokenType(typ types.Type) *types.Named {
	named := asNamed(typ)
	if named == nil || named.Obj().Pkg().Path() != "go/token" {
		return nil
	}
	return named
}

func asAstType(typ types.Type) *types.Named {
	typ = skipPtr(typ)
	named := asNamed(typ)
	if named == nil || named.Obj().Pkg().Path() != "go/ast" {
		return nil
	}
	return named
}

func asBasic(typ types.Type) *types.Basic {
	basic, _ := typ.(*types.Basic)
	return basic
}

func isSlice(obj types.Type) bool {
	_, ok := obj.(*types.Slice)
	return ok
}

func asSlice(obj types.Type) *types.Slice {
	slice, _ := obj.(*types.Slice)
	return slice
}

func asInterface(typ types.Type) *types.Interface {
	iface, _ := asNamed(typ).Underlying().(*types.Interface)
	return iface
}

func title(s string) string {
	if s != "" {
		return strings.ToUpper(s[:1]) + s[1:]
	}
	return ""
}

func zeroValue(typ types.Type) string {
	switch typ := typ.(type) {
	case *types.Basic:
		switch typ.Kind() {
		case types.String:
			return `""`
		case types.Bool:
			return `false`
		default:
			return `0`
		}
	case *types.Named:
		return zeroValue(typ.Underlying())
	case *types.Struct:
		return fmt.Sprintf("%s{}", typ.String())
	case *types.Interface:
		return `nil`
	case *types.Pointer:
		return `nil`
	case *types.Slice:
		return `nil`
	}
	return `"❌zero"`
}
