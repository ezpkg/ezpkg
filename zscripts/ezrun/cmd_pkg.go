package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/urfave/cli/v2"
	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/packages"

	"ezpkg.io/-/zscripts/-/script"
	"ezpkg.io/bytez"
	"ezpkg.io/errorz"
	"ezpkg.io/mapz"
)

const defaultPkgMode = packages.NeedName |
	packages.NeedFiles |
	packages.NeedCompiledGoFiles |
	packages.NeedImports |
	packages.NeedDeps |
	// packages.NeedExportFile |
	packages.NeedTypes |
	packages.NeedSyntax |
	packages.NeedTypesInfo |
	packages.NeedTypesSizes |
	packages.NeedModule |
	packages.NeedEmbedFiles |
	packages.NeedEmbedPatterns

type cmdPkg struct {
	fileSet *token.FileSet
	pkgInfo map[string]*PkgInfo
}

func (c *cmdPkg) Run(cx *cli.Context) error {
	c.pkgInfo = map[string]*PkgInfo{}

	args := script.WrapArgs(cx)
	switch {
	case cx.Bool("all"):
		args.MustEmpty()
		pkgs := listAllPkgs()
		c.generateCode(pkgs)

	default:
		pkgs := args.MustConsumeRemain(1, "NAME")
		c.generateCode(pkgs)
	}
	return nil
}

func (c *cmdPkg) generateCode(pkgs []string) {
	// clean all target dirs
	mustDir(env.TargetDir)
	for _, pkg := range pkgs {
		targetPkgDir := filepath.Join(env.TargetDir, pkg)
		items, err := os.ReadDir(targetPkgDir)
		switch {
		case err == nil:
			fmt.Printf("cleaning %v\n", targetPkgDir)
			for _, d := range items {
				if d.Name() == ".git" {
					continue
				}
				errorz.MustZ(os.RemoveAll(filepath.Join(targetPkgDir, d.Name())))
			}
		case errors.Is(err, os.ErrNotExist):
			errorz.MustZ(os.MkdirAll(targetPkgDir, 0755))
		default:
			errorz.MustZ(err)
		}
	}

	// parse packages
	var err error
	c.fileSet = token.NewFileSet()
	for _, pkg := range pkgs {
		fmt.Printf("loading %v\n", pkg)
		pkgInfo, err0 := c.processPackage(pkg)
		c.pkgInfo[pkg] = pkgInfo
		errorz.Append(&err, err0)
	}
	errorz.MustZ(err)

	// copy all files to target dirs, except tests and *_test.go
	for _, pkg := range pkgs {
		fmt.Printf("writing %v\n", pkg)
		pkgInfo := c.pkgInfo[pkg]
		pkgDir := filepath.Join(env.EzpkgDir, pkg)
		targetPkgDir := filepath.Join(env.TargetDir, pkg)

		for _, file := range pkgInfo.CompiledGoFiles {
			file = errorz.Must(filepath.Rel(pkgDir, file))
			copyFile(pkgDir, targetPkgDir, file)
		}
		// go.mod & go.local.mod
		errorz.MustZ(os.WriteFile(filepath.Join(targetPkgDir, "go.mod"), pkgInfo.goModPublish, 0644))
		errorz.MustZ(os.WriteFile(filepath.Join(targetPkgDir, "go.local.mod"), pkgInfo.goModLocal, 0644))
	}
	fmt.Printf("\n✅ DONE!\n")
}

func (c *cmdPkg) processPackage(pkgName string) (pkgInfo *PkgInfo, err error) {
	pkgDir := filepath.Join(env.EzpkgDir, pkgName)
	config := &packages.Config{
		Mode: defaultPkgMode,
		Fset: c.fileSet,
		Dir:  pkgDir,
	}
	_pkgs := errorz.Must(packages.Load(config, "./..."))
	for _, pkg := range _pkgs {
		errorz.Validatef(&err, strings.HasPrefix(pkg.PkgPath, "ezpkg.io/"), "package path must start with ezpkg.io: %v", pkg.PkgPath)
		if pkg.Name == pkgName {
			if pkgInfo != nil {
				panic(fmt.Sprintf("duplicated package name %q", pkg.Name))
			}
			pkgInfo = &PkgInfo{
				Package: pkg,
				pkgDir:  pkgDir,
			}
			c.processGoMod(pkgInfo, pkgDir)
		}
	}
	return pkgInfo, err
}

func (c *cmdPkg) processGoMod(pkgInfo *PkgInfo, pkgDir string) {
	path := filepath.Join(pkgDir, "go.mod")
	data := errorz.Must(os.ReadFile(path))
	pkgInfo.goMod = errorz.Must(modfile.Parse(path, data, nil))

	re := regexp.MustCompile(`\ngo [\d.]+\n`)
	idx := re.FindIndex(data)[1]

	var ezpkgImports []string
	for _, importPath := range mapz.SortedKeys(pkgInfo.Imports) {
		if strings.HasPrefix(importPath, "ezpkg.io/") {
			ezpkgImports = append(ezpkgImports, importPath)
		}
	}

	outputGoMod := func(isLocal bool) []byte {
		var b bytez.Buffer
		b.WriteZ(bytes.TrimSpace(data[:idx]))
		b.WriteStringZ("\n")
		if len(ezpkgImports) > 0 {
			b.Printf("\nrequire (\n")
			for _, importPath := range ezpkgImports {
				b.Printf("\t%s v%s\n", importPath, env.Info.Version)
			}
			b.Printf(")\n")
			if isLocal {
				b.Printf("replace (\n")
				for _, importPath := range ezpkgImports {
					b.Printf("\t%s => ../%s\n", importPath, filepath.Base(importPath))
				}
				b.Printf(")\n\n")
			}
		}
		data0 := bytes.TrimSpace(data[idx:])
		if len(data0) > 0 {
			b.WriteStringZ("\n")
			b.WriteZ(data0)
			b.WriteStringZ("\n")
		}
		return b.Bytes()
	}
	pkgInfo.goModPublish = outputGoMod(false)
	pkgInfo.goModLocal = outputGoMod(true)
}

func copyFile(srcDir, dstDir, file string) {
	srcPath := filepath.Join(srcDir, file)
	dstPath := filepath.Join(dstDir, file)
	data := errorz.Must(os.ReadFile(srcPath))
	errorz.MustZ(os.MkdirAll(filepath.Dir(dstPath), 0755))
	errorz.MustZ(os.WriteFile(dstPath, data, 0644))
}
