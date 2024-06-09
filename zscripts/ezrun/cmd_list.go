package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"

	"ezpkg.io/-/zscripts/-/script"
	"ezpkg.io/errorz"
	"ezpkg.io/unsafez"
)

type cmdList struct {
}

func (c *cmdList) Run(cx *cli.Context) error {
	pkgs := listAllPkgs()
	if len(pkgs) == 0 {
		script.Exitf("no packages found")
	}
	for _, pkg := range pkgs {
		fmt.Println(pkg)
	}
	return nil
}

func listAllPkgs() (pkgs []string) {
	goWork := unsafez.BytesToString(errorz.Must(os.ReadFile(filepath.Join(env.EzpkgDir, "go.work"))))
	for _, line := range strings.Split(goWork, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasSuffix(line, "z") {
			pkgs = append(pkgs, line)
		}
	}
	return pkgs
}
