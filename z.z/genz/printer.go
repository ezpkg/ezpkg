package genz

import (
	"bytes"
	"fmt"
	"go/types"
	"io"
	"regexp"
	"strconv"
)

type Printer interface {
	PkgPath() string // note: may be empty
	FilePath() string
	Import(name, path string)
	Qualifier(pkg *types.Package) string
	TypeString(types.Type) string
	Printf(msg string, args ...any)
	Bytes() []byte

	GetPkgPathByImportAlias(string) string

	io.WriteCloser
}

var _ Printer = &printer{}

type printer struct {
	engine   *engine
	plugin   *pluginStruct
	pkg      *types.Package // note that package can be nil
	pkgName  string
	filePath string
	closed   bool
	buf      *bytes.Buffer

	aliasByPkgPath map[string]string
	pkgPathByAlias map[string]string
}

func newPrinter(engine *engine, plugin *pluginStruct, pkg *types.Package, pkgName string, filePath string) *printer {
	if pkg != nil {
		pkgName = pkg.Name()
	}
	return &printer{
		engine:   engine,
		plugin:   plugin,
		pkg:      pkg,
		pkgName:  pkgName,
		filePath: filePath,
		buf:      engine.bufPool.Get().(*bytes.Buffer),

		aliasByPkgPath: make(map[string]string),
		pkgPathByAlias: make(map[string]string),
	}
}

func (p *printer) FilePath() string {
	return p.filePath
}

func (p *printer) PkgPath() string {
	return GetPkgPath(p.pkg)
}

func (p *printer) Printf(msg string, args ...any) {
	_, err := fmt.Fprintf(p, msg, args...)
	must(err)
}

func (p *printer) Write(data []byte) (n int, err error) {
	if p.closed {
		panic("already closed")
	}
	return p.buf.Write(data)
}

func (p *printer) Bytes() []byte {
	return p.buf.Bytes()
}

func (p *printer) Close() (_err error) {
	if p.closed {
		return nil
	}
	p.closed = true
	defer func() {
		p.buf.Reset()
		p.engine.bufPool.Put(p.buf)
	}()

	w, err := p.engine.writeFile(p.filePath)
	if err != nil {
		return err
	}
	defer func() {
		err2 := w.Close()
		if _err == nil {
			_err = err2
		}
	}()

	fprintf := func(format string, args ...any) {
		if _err != nil {
			return
		}
		_, err2 := fmt.Fprintf(w, format, args...)
		if err2 != nil {
			_err = err2
			return
		}
	}
	fprintf("//go:build !genz\n")
	fprintf("// Code generated by genz %v. DO NOT EDIT.\n\n", p.plugin.name)
	fprintf("package %v\n\n", p.pkgName)
	fprintf("import (\n")
	for path, alias := range p.aliasByPkgPath {
		fprintf("%v %q\n", alias, path)
	}
	fprintf(")\n\n")
	fprintf("%s", cleanCode(p.buf.Bytes()))
	return
}

func (p *printer) Import(name, path string) {
	if p.pkg != nil && p.pkg.Path() == path {
		return
	}
	if _, ok := p.aliasByPkgPath[path]; ok {
		return
	}
	if name == "" {
		p.aliasByPkgPath[path] = ""
		return
	}
	alias := name
	if _, ok := p.pkgPathByAlias[alias]; ok {
		for c := 1; ; c++ {
			alias = name + strconv.Itoa(c)
			if _, exist := p.pkgPathByAlias[alias]; !exist {
				break
			}
		}
	}
	p.pkgPathByAlias[alias] = path
	p.aliasByPkgPath[path] = alias
}

func (p *printer) GetPkgPathByImportAlias(alias string) string {
	return p.pkgPathByAlias[alias]
}

func (p *printer) TypeString(typ types.Type) string {
	return types.TypeString(typ, p.Qualifier)
}

func (p *printer) Qualifier(pkg *types.Package) string {
	if pkg == p.pkg {
		return ""
	}
	alias := pkg.Name()
	if p.plugin.qualifier != nil {
		alias = p.plugin.qualifier(pkg)
	}
	pkgPath := pkg.Path()
	p.Import(alias, pkgPath)
	return p.aliasByPkgPath[pkgPath]
}

// Clean pattern: {\n\n | \n\n}
var reClean1 = regexp.MustCompile(`\{\s*\n\s*\n`)
var reClean2 = regexp.MustCompile(`\n\s*\n\s*\}`)

func cleanCode(data []byte) []byte {
	data = reClean1.ReplaceAll(data, []byte("{\n"))
	data = reClean2.ReplaceAll(data, []byte("\n}"))
	return data
}
