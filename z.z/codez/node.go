package codez

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/packages"

	"ezpkg.io/slicez"
)

type _NodeReplaceFunc func(parent ast.Node, idx int, new ast.Node) error

type _NodeI interface{ Pos() token.Pos }
type _NodeW interface{ Unwrap() ast.Node }

type NodeX struct {
	ast.Node
	parent   ast.Node
	nodeIdx  int
	replacer _NodeReplaceFunc
}

func newNodeX(node ast.Node, parent ast.Node, nodeIdx int, replacer _NodeReplaceFunc) *NodeX {
	return &NodeX{
		Node:     node,
		parent:   parent,
		nodeIdx:  nodeIdx,
		replacer: replacer,
	}
}

func unwrapNode(node _NodeI) ast.Node {
	switch x := node.(type) {
	case *NodeX:
		return x.Node
	case ast.Node:
		return x
	}
	panic(fmt.Sprintf("unexpected node type: %T", node))
}

// intentionally to make NodeX not implement ast.Node
func (n *NodeX) End() {}

func (n *NodeX) Unwrap() ast.Node { return n.Node }

func (n *NodeX) Clone() *NodeX {
	return newNodeX(n.Node, n.parent, n.nodeIdx, n.replacer)
}

func (n *NodeX) ReplaceBy(new ast.Node) error {
	return n.replacer(n.parent, n.nodeIdx, new)
}

type FileX struct {
	*ast.File
}

func newFileX(file *ast.File) *FileX {
	return &FileX{
		File: file,
	}
}

// intentionally to make FileX not implement ast.Node
func (n *FileX) End() {}

func (f *FileX) Unwrap() ast.Node { return f.File }

func (f *FileX) GetImport(pkgPath string) *ast.ImportSpec {
	for _, imp := range f.Imports {
		if imp.Path.Value == fmt.Sprintf("%q", pkgPath) {
			return imp
		}
	}
	return nil
}

// Import adds an import to the file if it doesn't exist yet, and return the qualifier.
func (f *FileX) Import(alias string, pkg *packages.Package) (qualifier string) {
	getQualifier := func(imp *ast.ImportSpec) string {
		switch {
		case imp.Name == nil, imp.Name.Name == "":
			return pkg.Name
		case imp.Name.Name == ".": // dot import
			return ""
		default:
			return imp.Name.Name
		}
	}
	qualExists := func(qual string) bool {
		return slicez.ExistsFunc(f.Imports, func(imp *ast.ImportSpec) bool {
			return getQualifier(imp) == qual
		})
	}

	// if the package is already imported, return the qualifier
	if imp := f.GetImport(pkg.PkgPath); imp != nil {
		return getQualifier(imp)
	}

	// prepare the import spec
	imp := &ast.ImportSpec{
		Path: BasicString(pkg.PkgPath),
	}
	switch alias {
	case "":
		slicez.AppendTo(&f.Imports, imp)
		return pkg.Name

	case ".": // dot import
		imp.Name = NewIdent(".")
		slicez.AppendTo(&f.Imports, imp)
		return ""
	}

	// if the alias is already used, find a new one
	// TODO: handle conflicting default import alias
	for idx := 0; qualExists(alias); idx++ {
		alias = fmt.Sprintf("%s%d", alias, idx)
	}
	imp.Name = NewIdent(alias)
	slicez.AppendTo(&f.Imports, imp)

	return alias
}
