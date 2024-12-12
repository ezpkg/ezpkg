package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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
	Name      string    `json:"name"`
	Color     Color     `json:"color"`
	Colors    []Color   `json:"-"`
	MapColors MapColors `json:"colors"`
}

func (c *cmdList) Run(cx *cli.Context) error {
	pkgs := listAllPkgs()
	if len(pkgs) == 0 {
		script.Exitf("no packages found")
	}
	defer func() { // verify sorted
		for i := 0; i < len(pkgs)-1; i++ {
			if pkgs[i] > pkgs[i+1] {
				panic(fmt.Sprintf("invalid order of package: %v -> %v", pkgs[i], pkgs[i+1]))
			}
		}
	}()

	args := script.WrapArgs(cx)
	if !cx.Bool("x") {
		for _, pkg := range pkgs {
			fmt.Println(pkg)
		}
		return nil
	}
	ps := c.addExtraInfo(pkgs)
	fmt.Println(c.formatPkgs(ps))
	file := args.Consume()
	if file == "" {
		return nil
	}
	if !strings.HasSuffix(file, ".json") {
		script.Exitf("output file: only .json is supported")
	}

	mp := make(map[string]MapColors, len(ps))
	for _, pkg := range ps {
		mp[pkg.Name] = pkg.MapColors
	}
	data := errorz.Must(json.MarshalIndent(mp, "", "    "))
	return os.WriteFile(file, data, 0644)
}

func (c *cmdList) addExtraInfo(pkgs []string) []*cmdListPkgInfo {
	errorzIdx := slices.Index(pkgs, "errorz") // errorz is red
	if errorzIdx < 0 {
		panic("errorz not found")
	}

	gColor := GradientTable_Tailwind()
	N, ps := len(pkgs), make([]*cmdListPkgInfo, len(pkgs))
	for i := range pkgs {
		ps[i] = &cmdListPkgInfo{Name: pkgs[i]}
		pos := float64((i+N-errorzIdx)%N) / float64(N)
		ps[i].Color = gColor.GetInterpolatedColorFor(500, pos)
		ps[i].Colors = gColor.GetInterpolatedPaletteFor(pos)
		ps[i].MapColors = mapColors(gColor.Codes, ps[i].Colors)
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
		b.Printf("% 12s %s ", pkg.Name, pkg.Color.Hex())
		printColor(pkg.Color, pkg.Colors[2])
		b.Print(" ")
		for i, color := range pkg.Colors {
			textColor := typez.If(i > 4, pkg.Colors[2], pkg.Colors[len(pkg.Colors)-3])
			printColor(color, textColor)
		}
		b.Println()
	}
	return b.String()
}

func mapColors(codes []int, colors []Color) MapColors {
	if len(codes) != len(colors) {
		panic(fmt.Sprintf("wrong number of colors (expected %v)", len(codes)))
	}
	mp := make(MapColors, len(colors))
	for i, code := range codes {
		mp[i].Code = code
		mp[i].Color = colors[i]
	}
	return mp
}

func listAllPkgs() (pkgs []string) {
	rePkg := regexp.MustCompile(`^\w\w+\.\w\w+$`) // iter.json, not z.z

	goWork := unsafez.BytesToString(errorz.Must(os.ReadFile(filepath.Join(env.EzpkgDir, "go.work"))))
	for _, line := range strings.Split(goWork, "\n") {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "/") {
			continue
		}
		switch {
		case strings.HasSuffix(line, "z"):
			pkgs = append(pkgs, line)
		case rePkg.MatchString(line):
			pkgs = append(pkgs, line)
		}

	}
	return pkgs
}
