package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/urfave/cli/v2"

	"ezpkg.io/-/zscripts/-/script"
	"ezpkg.io/colorz"
	"ezpkg.io/errorz"
	"ezpkg.io/stringz"
	"ezpkg.io/typez"
	"ezpkg.io/unsafez"
)

type cmdList struct {
}

type cmdListPkgInfo struct {
	pkg    string
	color  Color
	colors []Color
}

func (c *cmdList) Run(cx *cli.Context) error {
	pkgs := listAllPkgs()
	if len(pkgs) == 0 {
		script.Exitf("no packages found")
	}
	if cx.Bool("x") {
		ps := c.addExtraInfo(pkgs)
		fmt.Println(c.formatPkgs(ps))
	} else {
		for _, pkg := range pkgs {
			fmt.Println(pkg)
		}
	}
	// verify sorted
	for i := 0; i < len(pkgs)-1; i++ {
		if pkgs[i] > pkgs[i+1] {
			panic(fmt.Sprintf("invalid order of package: %v -> %v", pkgs[i], pkgs[i+1]))
		}
	}
	return nil
}

func (c *cmdList) addExtraInfo(pkgs []string) []*cmdListPkgInfo {
	errorzIdx := slices.Index(pkgs, "errorz") // errorz is red
	if errorzIdx < 0 {
		panic("errorz not found")
	}

	gColor := GradientTable_Tailwind()
	N, ps := len(pkgs), make([]*cmdListPkgInfo, len(pkgs))
	for i := range pkgs {
		ps[i] = &cmdListPkgInfo{pkg: pkgs[i]}
		pos := float64((i+N-errorzIdx)%N) / float64(N)
		ps[i].color = gColor.GetInterpolatedColorFor(500, pos)
		ps[i].colors = gColor.GetInterpolatedPaletteFor(pos)
	}
	return ps
}

func (c *cmdList) formatPkgs(pkgs []*cmdListPkgInfo) string {
	var b stringz.Builder
	printColor := func(color, text Color) {
		R, G, B := color.RGB255()
		b.Printf("\x1b[48;2;%d;%d;%dm", R, G, B)
		R, G, B = text.RGB255()
		b.Printf("\x1b[38;2;%d;%d;%dm", R, G, B)
		b.Printf("  ⏺  %s", colorz.Reset)
	}
	for _, pkg := range pkgs {
		b.Printf("% 12s %s ", pkg.pkg, pkg.color.Hex())
		printColor(pkg.color, pkg.colors[2])
		b.Print(" ")
		for i, color := range pkg.colors {
			textColor := typez.If(i > 4, pkg.colors[2], pkg.colors[len(pkg.colors)-3])
			printColor(color, textColor)
		}
		b.Println()
	}
	return b.String()
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
