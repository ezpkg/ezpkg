package ggen

import (
	"fmt"
	"go/ast"
	"go/types"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"golang.org/x/tools/go/packages"

	"ezpkg.io/errorz"
	_ "ezpkg.io/genz/builtin"
	"ezpkg.io/logz"
)

const defaultGeneratedFileNameTpl = "zz_generated.%v.go"
const defaultBufSize = 1024 * 4
const startDirectiveStr0 = "//+"   //   //+directive
const startDirectiveStr1 = "// +"  //   //+directive
const startDirectiveStr2 = "//go:" //   //go:build ggen

var reAlphabet = regexp.MustCompile(`[a-z]`)
var reCommand = regexp.MustCompile(`^[a-z]([a-z0-9.:-]*[a-z0-9])?$`)

type Logger = logz.Logger
type zLogger struct{ Logger }

func FilterByCommand(command string) CommandFilter {
	return CommandFilter(command)
}

type CommandFilter string

func (cmd CommandFilter) Filter(ng FilterEngine) error {
	for _, p := range ng.ParsingPackages() {
		if cmd.Include(p.Directives) {
			p.Include()
		}
	}
	return nil
}

func (cmd CommandFilter) FilterAll(ng FilterEngine) error {
	for _, p := range ng.ParsingPackages() {
		if cmd.Include(p.Directives) {
			p.Include()
		} else if cmd.Include(p.InlineDirectives) {
			p.Include()
		}
	}
	return nil
}

func (cmd CommandFilter) Include(ds Directives) bool {
	for _, d := range ds {
		if d.Cmd == string(cmd) ||
			strings.HasPrefix(d.Cmd, string(cmd)) && d.Cmd[len(cmd)] == ':' {
			return true
		}
	}
	return false
}

func GetPkgDir(pkg *packages.Package) string {
	if pkg != nil && len(pkg.CompiledGoFiles) > 0 {
		return filepath.Dir(pkg.CompiledGoFiles[0])
	}
	return ""
}

func GetPkgPath(pkg *types.Package) string {
	if pkg != nil {
		return pkg.Path()
	}
	return ""
}

func GetPkgPathOfType(typ types.Type) string {
	for {
		ptr, ok := typ.(*types.Pointer)
		if !ok {
			break
		}
		typ = ptr.Elem()
	}
	named, ok := typ.(*types.Named)
	if !ok {
		return ""
	}
	pkg := named.Obj().Pkg()
	if pkg != nil {
		return pkg.Path()
	}
	return ""
}

func defaultFileNameGenerator(tpl string) func(GenerateFileNameInput) string {
	return func(input GenerateFileNameInput) string {
		return fmt.Sprintf(tpl, strings.ReplaceAll(input.PluginName, "-", "_"))
	}
}

var ggenPath = reflect.TypeOf((*Engine)(nil)).Elem().PkgPath()
var builtinPath = filepath.Dir(ggenPath) + "/builtin"

func parseBuiltinTypes(pkg *packages.Package) map[string]types.Type {
	if pkg.PkgPath != builtinPath {
		panic(fmt.Sprintf("unexpected path %v", pkg.PkgPath))
	}
	m := map[string]types.Type{}
	s := pkg.Types.Scope()
	for _, name := range s.Names() {
		if !strings.HasPrefix(name, "_") {
			continue
		}
		typ := s.Lookup(name).Type()
		m[typ.String()] = typ
	}
	return m
}

func getPackageDir(pkg *packages.Package) string {
	if len(pkg.GoFiles) > 0 {
		return filepath.Dir(pkg.GoFiles[0])
	}
	return ""
}

func hasStartDirective(line string) bool {
	return strings.HasPrefix(line, startDirectiveStr0) ||
		strings.HasPrefix(line, startDirectiveStr1) ||
		strings.HasPrefix(line, startDirectiveStr2)
}

// processDoc splits directive and text comment
func processDoc(doc, cmt *ast.CommentGroup) (Comment, error) {
	if doc == nil {
		return Comment{Comment: cmt}, nil
	}

	directives := make([]Directive, 0, 4)
	for _, line := range doc.List {
		if !hasStartDirective(line.Text) {
			continue
		}

		// remove "// " but keep "+"
		text := strings.TrimSpace(strings.TrimPrefix(line.Text, "//"))
		directive, err := ParseDirective(text)
		if err != nil {
			return Comment{}, err
		}
		directives = append(directives, directive)
	}

	comment := Comment{
		Doc:        doc,
		Comment:    cmt,
		Directives: directives,
	}
	return comment, nil
}

func processDocText(doc *ast.CommentGroup) string {
	if doc == nil {
		return ""
	}
	processedDoc := make([]*ast.Comment, 0, len(doc.List))
	for _, line := range doc.List {
		if hasStartDirective(line.Text) {
			processedDoc = append(processedDoc, line)
			continue
		}
	}
	return (&ast.CommentGroup{List: processedDoc}).Text()
}

// ParseDirectiveFromFile reads from file and returns the parsed directives.
func ParseDirectiveFromFile(filename string) (directives, inlineDirective []Directive, err error) {
	body, err := os.ReadFile(filename)
	if err != nil {
		return
	}
	return ParseDirectiveFromBody(body)
}

// ParseDirectiveFromBody reads directives from body and returns the parsed directives.
func ParseDirectiveFromBody(body []byte) (directives, inlineDirective []Directive, err error) {
	errs := parseDirectivesFromBody(body, &directives, &inlineDirective)
	err = errorz.Wrap(errorz.Combine(errs...), "can not parse directive")
	return
}

// ParseDirective parses directives from a single line.
func ParseDirective(line string) (result Directive, _ error) {
	line = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "//"))
	if line == "go:build" || strings.HasPrefix(line, "go:build ") {
		return parseGoBuildDirective(line)
	}

	result, err := parsePlusDirective(line)
	if err != nil {
		return result, errorz.Wrapf(err, "%v (%v)", err, line)
	}
	return result, nil
}

func parseGoBuildDirective(line string) (Directive, error) {
	arg := strings.TrimPrefix(line, "go:build")
	arg = strings.TrimSpace(arg)
	directive := Directive{
		Raw: line,
		Cmd: "go:build",
		Arg: arg,
	}
	return directive, nil
}

func parsePlusDirective(line string) (result Directive, err error) {
	result.Raw = line
	line = strings.TrimPrefix(line, "+")
	idx := strings.IndexByte(line, ' ') //   //+name arg
	if idx >= 0 {
		result.Cmd = line[:idx]
		result.Arg = strings.TrimSpace(line[idx+1:])
	} else {
		result.Cmd = line
	}
	if reAlphabet.MatchString(result.Cmd) && !reCommand.MatchString(result.Cmd) {
		return result, errorz.NoStack().New("invalid directive")
	}
	return result, nil
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
