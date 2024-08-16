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
	"ezpkg.io/slicez"
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

	{
		_, astExprI := getIface(pkgGoAst, "Expr")
		p := errorz.Must(ng.GenerateFile("codez", pkgDir+"/matchers.expr.go"))
		p.Import("ast", "go/ast")

		defer func() { errorz.MustZ(p.Close()) }()
		pr := p.Printf

		for _, x := range allObjs.Implements(astExprI).Structs() {
			pr("// %s\n", x.Name())
			pr("type %sMatcher struct {\n", x.Name())
			pr("\t_ *%s\n\n", p.TypeString(x.Type()))
			st := mustStruct(x.Type())
			for i := 0; i < st.NumFields(); i++ {
				field := st.Field(i)
				if typ := asTokenType(field.Type()); typ != nil {
					pr("\t%s %s\n", field.Name(), p.TypeString(field.Type()))
					continue
				}
				if typ := asAstType(field.Type()); typ != nil {
					pr("\t%s %sMatcher\n", field.Name(), typ.Obj().Name())
					continue
				}
				if slice := asSlice(field.Type()); slice != nil {
					typ := asAstType(slice.Elem())
					if typ == nil {
						panic(fmt.Sprintf("unsupported slice type %v", field.Type()))
					}
					pr("\t%s %sListMatcher\n", field.Name(), typ.Obj().Name())
					continue
				}
				if basic := asBasic(field.Type()); basic != nil {
					pr("\t%s %sMatcher\n", field.Name(), title(field.Type().String()))
					continue
				}
				pr("\t%s %s // ❌\n", field.Name(), p.TypeString(field.Type()))
			}
			pr("}\n\n")
		}
	}
	{
		_, astStmtI := getIface(pkgGoAst, "Stmt")
		p := errorz.Must(ng.GenerateFile("codez", pkgDir+"/matchers.stmt.go"))
		p.Import("ast", "go/ast")

		defer func() { errorz.MustZ(p.Close()) }()
		pr := p.Printf

		for _, x := range allObjs.Implements(astStmtI).Structs() {
			pr("// %s\n", x.Name())
			pr("type %sMatcher struct {\n", x.Name())
			pr("\t_ *%s\n\n", p.TypeString(x.Type()))
			st := mustStruct(x.Type())
			for i := 0; i < st.NumFields(); i++ {
				field := st.Field(i)
				if typ := asTokenType(field.Type()); typ != nil {
					pr("\t%s %s\n", field.Name(), p.TypeString(field.Type()))
					continue
				}
				if typ := asAstType(field.Type()); typ != nil {
					pr("\t%s %sMatcher\n", field.Name(), typ.Obj().Name())
					continue
				}
				if slice := asSlice(field.Type()); slice != nil {
					typ := asAstType(slice.Elem())
					if typ == nil {
						panic(fmt.Sprintf("unsupported slice type %v", field.Type()))
					}
					pr("\t%s %sListMatcher\n", field.Name(), typ.Obj().Name())
					continue
				}
				if basic := asBasic(field.Type()); basic != nil {
					pr("\t%s %sMatcher\n", field.Name(), title(field.Type().String()))
					continue
				}
				pr("\t%s %s // ❌\n", field.Name(), p.TypeString(field.Type()))
			}
			pr("}\n\n")
		}
	}
	{
		_, astDeclI := getIface(pkgGoAst, "Decl")
		p := errorz.Must(ng.GenerateFile("codez", pkgDir+"/matchers.decl.go"))
		p.Import("ast", "go/ast")

		defer func() { errorz.MustZ(p.Close()) }()
		pr := p.Printf

		for _, x := range allObjs.Implements(astDeclI).Structs() {
			pr("// %s\n", x.Name())
			pr("type %sMatcher struct {\n", x.Name())
			pr("\t_ *%s\n\n", p.TypeString(x.Type()))
			st := mustStruct(x.Type())
			for i := 0; i < st.NumFields(); i++ {
				field := st.Field(i)
				if typ := asTokenType(field.Type()); typ != nil {
					pr("\t%s %s\n", field.Name(), p.TypeString(field.Type()))
					continue
				}
				if typ := asAstType(field.Type()); typ != nil {
					pr("\t%s %sMatcher\n", field.Name(), typ.Obj().Name())
					continue
				}
				if slice := asSlice(field.Type()); slice != nil {
					typ := asAstType(slice.Elem())
					if typ == nil {
						panic(fmt.Sprintf("unsupported slice type %v", field.Type()))
					}
					pr("\t%s %sListMatcher\n", field.Name(), typ.Obj().Name())
					continue
				}
				if basic := asBasic(field.Type()); basic != nil {
					pr("\t%s %sMatcher\n", field.Name(), title(field.Type().String()))
					continue
				}
				pr("\t%s %s // ❌\n", field.Name(), p.TypeString(field.Type()))
			}
			pr("}\n\n")
		}
	}
	return nil
}

type Objects []types.Object

func (objs Objects) Implements(iface *types.Interface) (out Objects) {
	for _, obj := range objs {
		ptr := types.NewPointer(obj.Type())
		if types.Implements(obj.Type(), iface) || types.Implements(ptr, iface) {
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

func asPtrNamed(typ types.Type) *types.Named {
	ptr, _ := typ.(*types.Pointer)
	return asNamed(ptr.Elem())
}

func asNamed(typ types.Type) *types.Named {
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
