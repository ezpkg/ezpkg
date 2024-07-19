package tests_test

import (
	"go/types"
	"io"
	"os/exec"
	"path/filepath"
	"testing"

	"ezpkg.io/genz/ggen"

	"github.com/stretchr/testify/require"
)

const testPath = "ezpkg.io/genz/tests"
const testPatterns = testPath + "/..."

type mockPlugin struct {
	ng ggen.Engine

	filter   func(ggen.FilterEngine) error
	generate func(ggen.Engine) error
	imAlias  func(string, string) string
}

func (m *mockPlugin) Name() string    { return "mock" }
func (m *mockPlugin) Command() string { return "gen:mock" }

func (m *mockPlugin) Filter(ng ggen.FilterEngine) error {
	if m.filter != nil {
		return m.filter(ng)
	}
	for _, p := range ng.ParsingPackages() {
		p.Include()
	}
	return nil
}

func (m *mockPlugin) Generate(ng ggen.Engine) error {
	m.ng = ng
	if m.generate != nil {
		return m.generate(ng)
	}
	return nil
}

func (m *mockPlugin) ImportAlias(pkgPath, importPath string) string {
	if m.imAlias != nil {
		return m.imAlias(pkgPath, importPath)
	}
	return ""
}

var registered = false
var mock = &mockPlugin{}

func reset() {
	*mock = mockPlugin{} // reset the plugin
	if registered {
		return
	}
	registered = true
}

func TestObjects(t *testing.T) {
	reset()
	cfg := ggen.Config{}
	cfg.RegisterPlugin(mock)
	err := ggen.Start(cfg, testPatterns)
	require.NoError(t, err)

	ng := mock.ng
	pkg := ng.GetPackageByPath(testPath + "/one")
	require.NotNil(t, pkg)

	objects := ng.GetObjectsByScope(pkg.Types.Scope())
	require.Len(t, objects, 2)
	require.Equal(t, "A", objects[0].Name())
	require.Equal(t, "B", objects[1].Name())

	{
		directives := ng.GetDirectives(objects[0])
		require.Len(t, directives, 1)
		require.Equal(t, "ggen:a", directives[0].Cmd)
	}
	{
		directives := ng.GetDirectives(objects[1])
		require.Len(t, directives, 1)
		require.Equal(t, "ggen:b", directives[0].Cmd)
	}
	{
		objA := objects[0]
		cmt := ng.GetComment(objA)
		require.Equal(t, "this is comment of A\n", cmt.Text())
		require.Len(t, cmt.Directives, 1)
		require.Equal(t, "ggen:a", cmt.Directives[0].Cmd)

		st, ok := objA.Type().Underlying().(*types.Struct)
		require.True(t, ok, "should be *types.Struct")
		zero := st.Field(0)
		one := st.Field(1)
		two := st.Field(2)
		thr := st.Field(3)
		require.Equal(t, "", ng.GetComment(zero).Text())
		require.Equal(t, "comment of One\n", ng.GetComment(one).Text())
		require.Equal(t, "", ng.GetComment(two).Text())
		require.Equal(t, "comment of Three\n", ng.GetComment(thr).Text())
	}
}

func TestGenerate(t *testing.T) {
	reset()
	var pkgs []*ggen.GeneratingPackage
	mock.generate = func(ng ggen.Engine) error {
		pkgs = ng.GeneratingPackages()
		for _, pkg := range pkgs {
			// skip package "two"
			if pkg.Package.PkgPath == testPath+"/two" {
				_ = pkg.GetPrinter()
				continue
			}

			p := pkg.GetPrinter()
			mustWrite(p, []byte(" ")) // write a single byte for triggering p.Close()
		}
		return nil
	}

	cfg := ggen.Config{}
	cfg.RegisterPlugin(mock)
	err := ggen.Start(cfg, testPatterns)
	require.NoError(t, err)

	output, err := exec.Command("sh", "-c", `find . | grep zz | sort`).
		CombinedOutput()
	require.NoError(t, err)

	expected := `
./one/one-and-a-half/zz_generated.mock.go
./one/zz_generated.mock.go
./zz_generated.mock.go
`[1:]
	require.Equal(t, expected, string(output))
}

func TestClean(t *testing.T) {
	reset()
	cfg := ggen.Config{CleanOnly: true}
	cfg.RegisterPlugin(mock)
	err := ggen.Start(cfg, testPatterns)
	require.NoError(t, err)

	output, err := exec.Command("sh", "-c", `find . | grep zz | sort`).
		CombinedOutput()
	require.NoError(t, err)
	require.Equal(t, "", string(output))
}

func TestInclude(t *testing.T) {
	reset()

	parentPath := filepath.Dir(testPath)
	mock.filter = func(ng ggen.FilterEngine) error {
		ng.IncludePackage(testPath + "/two")
		ng.IncludePackage(parentPath) // parentPath is outside of testPatterns
		return nil
	}

	cfg := ggen.Config{}
	cfg.RegisterPlugin(mock)
	err := ggen.Start(cfg, testPatterns)
	require.NoError(t, err)

	expecteds := []string{
		parentPath,
		testPath + "/two",
	}
	pkgs := mock.ng.GeneratingPackages()
	require.Len(t, pkgs, 2)
	for i, pkg := range pkgs {
		require.Equal(t, expecteds[i], pkg.PkgPath)
	}
}

func mustWrite(w io.Writer, p []byte) {
	if _, err := w.Write(p); err != nil {
		panic(err)
	}
}
