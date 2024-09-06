package codez

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"io"
	"os"

	"golang.org/x/tools/go/packages"

	"ezpkg.io/errorz"
	"ezpkg.io/slicez"
	"ezpkg.io/typez"
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
	px  *Packages
	pkg *PackageX
	tok *token.File
}

func (px *Packages) newFileX(file *ast.File, tok *token.File, pkg *PackageX) *FileX {
	return &FileX{
		File: file,
		px:   px,
		pkg:  pkg,
		tok:  tok,
	}
}

// intentionally to make FileX not implement ast.Node
func (n *FileX) End() {}

func (f *FileX) Unwrap() ast.Node { return f.File }

func (p *FileX) Path() string {
	return p.tok.Name()
}

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

// WriteFile writes the file to the input path, creating if necessary. If path is empty, it writes to the original file. If the file exists, it will be overridden.
func (f *FileX) WriteFile(path string, perm os.FileMode) error {
	var b bytes.Buffer
	if err := format.Node(&b, f.px.Fset, f); err != nil {
		return errorz.Wrapf(err, "failed to format file")
	}
	path = typez.Coalesce(path, f.Path())
	if path == "" {
		return errorz.Errorf("empty path")
	}
	if err := os.WriteFile(path, b.Bytes(), perm); err != nil {
		return errorz.Wrapf(err, "failed to write file")
	}
	return nil
}

// WriteTo writes the file to the writer.
func (f *FileX) WriteTo(w io.Writer) (int64, error) {
	var _ io.WriterTo = (*FileX)(nil)

	var b bytes.Buffer
	if err := format.Node(&b, f.px.Fset, f); err != nil {
		return 0, errorz.Wrap(err, "failed to format file")
	}
	n, err := w.Write(b.Bytes())
	return int64(n), errorz.Wrap(err, "failed to write file")
}
