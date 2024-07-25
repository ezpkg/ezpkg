package ggen

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

type Comment struct {
	Doc *ast.CommentGroup

	Comment *ast.CommentGroup

	Directives []Directive
}

func (c Comment) Text() string {
	return processDocText(c.Doc)
}

// Directive comment has one of following formats
//
//	// +foo:valid=required,optional
//	// +foo:valid=null +gen=foo
//	// +foo:pkg=sample,baz
//	// +foo:valid: 0 < $ && $ <= 10
//
// For example "// +foo:pkg=sample,baz" will be parsed as
//
//	Command: "foo:pkg"
//	Arg:     "sample,baz"
//
// Directives must start at the begin of a line, after "//" and a space (the
// same as "// +build"). Multiple directives can appear in one line.
//
// Directive ending with "=" can not have space in argument and can have
// multiple directives. Directive ending with ":" can have space in argument,
// therefore it will be parsed as a single directive.
type Directive struct {
	Raw string // +foo:pkg:foo this is a string
	Cmd string // foo:pkg
	Arg string // sample,baz

	Item Positioner // the item that the directive is attached to
}

func (d Directive) String() string {
	return d.Raw
}

func (d Directive) IsPackageLevel() bool {
	return d.Item == nil
}

// ParseArgs parse directive argument using the standard "flag" package format. Example:
//
// // +ggen:sample -name=Alice DoSomething
func (d Directive) GetArgs() ([]string, error) {
	if d.Arg == "" {
		return nil, nil
	}
	// TODO(iolivernguyen): handle escape
	// // +ggen:sample -name="Alice M"
	return strings.Split(d.Arg, " "), nil
}

type Directives []Directive

func (ds Directives) Get(cmd string) (Directive, bool) {
	for _, d := range ds {
		if d.Cmd == cmd {
			return d, true
		}
	}
	return Directive{}, false
}

func (ds Directives) GetArg(cmd string) string {
	for _, d := range ds {
		if d.Cmd == cmd {
			return d.Arg
		}
	}
	return ""
}

// FilterBy returns list of directives that have the given command.
// Examples of accepted directives with input "+ggen:sample" or "ggen:sample"
//
// // +ggen:sample
// // +ggen:sample:foo
// // +ggen:sample:foo argument
func (ds Directives) FilterBy(prefix string) Directives {
	if strings.HasPrefix(prefix, "+") {
		prefix = prefix[1:]
	}
	if !strings.HasSuffix(prefix, ":") {
		prefix = prefix + ":"
	}
	res := make([]Directive, 0, len(ds))
	for _, d := range ds {
		if d.Cmd == prefix[:len(prefix)-1] ||
			strings.HasPrefix(d.Cmd, prefix) {
			res = append(res, d)
		}
	}
	return res
}

type declaration struct {
	Pkg *packages.Package

	Comment Comment
}

type extendedInfo struct {
	// FileSet
	Fset *token.FileSet

	// Map from Ident to declaration
	Declarations map[*ast.Ident]*declaration

	// Map from token.Pos to Ident
	Positions map[token.Pos]*ast.Ident
}

func newExtendedInfo(fset *token.FileSet) *extendedInfo {
	return &extendedInfo{
		Fset:         fset,
		Declarations: make(map[*ast.Ident]*declaration),
		Positions:    make(map[token.Pos]*ast.Ident),
	}
}

func (x *extendedInfo) AddPackage(pkg *packages.Package) error {
	for _, file := range pkg.Syntax {
		err := x.addFile(pkg, file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (x *extendedInfo) addFile(pkg *packages.Package, file *ast.File) error {
	var genDoc *ast.CommentGroup
	processDocFunc := func(doc, cmt *ast.CommentGroup, fallback bool) *declaration {
		if fallback {
			if doc == nil {
				doc = genDoc
			}
		} else {
			genDoc = nil
		}
		comment, err := processDoc(doc, cmt)
		if err != nil {
			logger.Debug("error while processing doc", "err", err)
		}
		return &declaration{
			Pkg:     pkg,
			Comment: comment,
		}
	}

	mergeCmt := func(a, b *ast.CommentGroup) *ast.CommentGroup {
		switch {
		case a == nil:
			return b
		case b == nil:
			return a
		case a == b:
			return a
		default:
			logger.Warn("conflicting comments", nil, "a", a.Text(), "b", b.Text())
			a.List = append(a.List, b.List...)
			return a
		}
	}
	// setDecl may be called multiple times for the same ident, so it will attempt to merge the declarations
	setDecl := func(ident *ast.Ident, decl *declaration) {
		if x.Declarations[ident] == nil {
			x.Declarations[ident] = decl
			return
		}
		prev := x.Declarations[ident]
		if prev.Pkg != decl.Pkg {
			logger.Warn("conflicting declarations", nil, ident.Name, "prev", prev.Pkg.PkgPath, "new", decl.Pkg.PkgPath)
		}
		prev.Comment.Doc = mergeCmt(prev.Comment.Doc, decl.Comment.Doc)
		prev.Comment.Comment = mergeCmt(prev.Comment.Comment, decl.Comment.Comment)
	}

	positions := x.Positions
	ast.Inspect(file, func(node ast.Node) bool {
		switch node := node.(type) {
		case *ast.Ident:
			setDecl(node, &declaration{Pkg: pkg})
			positions[node.NamePos] = node

		case *ast.FuncDecl:
			ident := node.Name
			setDecl(ident, processDocFunc(node.Doc, nil, false))
			positions[ident.NamePos] = ident

		case *ast.GenDecl:
			// if the declaration has only 1 spec, we treat the doc
			// comment as the spec comment
			if len(node.Specs) == 1 {
				genDoc = node.Doc
			} else {
				genDoc = nil
			}

		case *ast.ImportSpec:
			if node.Name != nil {
				ident := node.Name
				setDecl(ident, processDocFunc(node.Doc, node.Comment, true))
				positions[ident.Pos()] = ident
			}

		case *ast.TypeSpec:
			ident := node.Name

			setDecl(ident, processDocFunc(node.Doc, node.Comment, true))
			positions[ident.NamePos] = ident

		case *ast.ValueSpec:
			for _, ident := range node.Names {
				setDecl(ident, processDocFunc(node.Doc, node.Comment, true))
				positions[ident.NamePos] = ident
			}

		case *ast.Field:
			for _, ident := range node.Names {
				setDecl(ident, processDocFunc(node.Doc, node.Comment, false))
				positions[ident.NamePos] = ident
			}
		}
		return true
	})
	return nil
}

func (x *extendedInfo) GetObject(ident *ast.Ident) types.Object {
	decl := x.Declarations[ident]
	if decl == nil {
		return nil
	}
	return decl.Pkg.TypesInfo.ObjectOf(ident)
}

func (x *extendedInfo) GetComment(ident *ast.Ident) Comment {
	decl := x.Declarations[ident]
	if decl == nil {
		return Comment{}
	}
	return decl.Comment
}
