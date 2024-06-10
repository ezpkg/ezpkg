package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/urfave/cli/v2"

	"ezpkg.io/-/zscripts/-/script"
	"ezpkg.io/errorz"
	"ezpkg.io/unsafez"
)

type cmdDoc struct {
}

type ReadmeArgs struct {
	pkgName string
	version string // 0.0.1
	pkgDesc string
	pkgDoc  string
}

func (c *cmdDoc) Run(cx *cli.Context) error {
	args := script.WrapArgs(cx)
	switch {
	case cx.Bool("all"):
		args.MustEmpty()
		pkgs := listAllPkgs()
		c.generateDoc(pkgs)

	default:
		pkgs := args.MustConsumeRemain(1, "NAME")
		c.generateDoc(pkgs)
	}
	return nil
}

func (c *cmdDoc) generateDoc(pkgs []string) {
	tpl := c.loadReadmeTpl()
	for _, pkg := range pkgs {
		fmt.Printf("👉 %s/README.md\n", pkg)
		args := c.loadPkgDoc(pkg)
		data := c.formatTpl(tpl, args)

		path := env.EzpkgDir + "/" + pkg + "/README.md"
		errorz.MustZ(os.WriteFile(path, data, 0644))
	}
	fmt.Printf("\n✅ DONE!\n")
}

func (c *cmdDoc) loadPkgDoc(pkg string) ReadmeArgs {
	path := env.EzpkgDir + "/" + pkg + "/DOC.md"
	raw := unsafez.BytesToString(errorz.Must(os.ReadFile(path)))
	raw = strings.TrimSpace(raw)
	parts := strings.SplitN(raw, "\n", 2)
	return ReadmeArgs{
		pkgName: pkg,
		version: env.Info.Version,
		pkgDesc: strings.TrimSpace(parts[0]),
		pkgDoc:  strings.TrimSpace(parts[1]),
	}
}

func (c *cmdDoc) loadReadmeTpl() []byte {
	path := env.ZscriptsDir + "/zpkg.tpl/README.md"
	return errorz.Must(os.ReadFile(path))
}

func (c *cmdDoc) formatTpl(tpl []byte, args ReadmeArgs) []byte {
	check := func(s string, name string) []byte {
		s = strings.TrimSpace(s)
		if s == "" {
			panic("empty " + name)
		}
		return []byte(s)
	}
	replaces := map[string][]byte{
		"{PKGNAME}":  check(args.pkgName, "pkgName"),
		"{VERSION}":  check(args.version, "version"),
		"{PKG_DESC}": check(args.pkgDesc, "pkgDesc"),
		"{PKG_DOC}":  check(args.pkgDoc, "pkgDoc"),
	}
	return regexp.MustCompile(`\{[A-Z_]+}`).ReplaceAllFunc(tpl, func(s []byte) []byte {
		out := replaces[unsafez.BytesToString(s)]
		if out == nil {
			panic("no replacement for " + string(s))
		}
		return out
	})
}
