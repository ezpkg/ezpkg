package genz

import (
	"go/types"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"
)

func TestBuiltinTypes(t *testing.T) {
	cfg := &packages.Config{Mode: packages.LoadAllSyntax}
	pkgs, err := packages.Load(cfg, builtinPath)
	require.NoError(t, err)

	pkg := pkgs[0]
	require.NotNil(t, pkg)

	m := parseBuiltinTypes(pkg)
	require.Equal(t, types.Int, m["int"].(*types.Basic).Kind())
	require.Equal(t, "Error", m["error"].Underlying().(*types.Interface).Method(0).Name())
}

func TestParseGoBuild(t *testing.T) {
	t.Run("go:build", func(t *testing.T) {
		directive, err := ParseDirective("go:build")
		require.NoError(t, err)
		require.Equal(t, Directive{
			Raw: "go:build",
			Cmd: "go:build",
			Arg: "",
		}, directive)
	})
	t.Run("go:build genz", func(t *testing.T) {
		directive, err := ParseDirective("go:build genz")
		require.NoError(t, err)
		require.Equal(t, Directive{
			Raw: "go:build genz",
			Cmd: "go:build",
			Arg: "genz",
		}, directive)
	})
}

func TestParseDirective(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		directive, err := parsePlusDirective("")
		require.NoError(t, err)
		require.Empty(t, directive)
	})
	t.Run("ignore", func(t *testing.T) {
		_, err := parsePlusDirective("+++")
		require.NoError(t, err)
	})
	t.Run("invalid", func(t *testing.T) {
		_, err := parsePlusDirective("+a+")
		require.Error(t, err)
	})
	t.Run("sample", func(t *testing.T) {
		directive, err := parsePlusDirective("+sample")
		require.NoError(t, err)
		require.Equal(t, Directive{
			Raw: "+sample",
			Cmd: "sample",
			Arg: "",
		}, directive)
	})
	t.Run("sample with args", func(t *testing.T) {
		directive, err := parsePlusDirective("+sample arg1,arg2")
		require.NoError(t, err)
		require.Equal(t, Directive{
			Raw: "+sample arg1,arg2",
			Cmd: "sample",
			Arg: "arg1,arg2",
		}, directive)
	})
	t.Run("sample with spaces in args", func(t *testing.T) {
		directive, err := parsePlusDirective("+sample -key1=arg1 -key2=arg2")
		require.NoError(t, err)
		require.Equal(t, Directive{
			Raw: "+sample -key1=arg1 -key2=arg2",
			Cmd: "sample",
			Arg: "-key1=arg1 -key2=arg2",
		}, directive)
	})
}
