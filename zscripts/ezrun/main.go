package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"

	"ezpkg.io/-/zscripts/-/script"
	"ezpkg.io/errorz"
)

//go:embed _usage.txt
var usageText string

var flags struct {
	Verbose bool
}

var env struct {
	RootDir     string
	EzpkgDir    string
	TargetDir   string
	ZscriptsDir string

	Info *EzpkgInfo
}

var flagAll = &cli.BoolFlag{Name: "all", Usage: "all packages"}
var flagX = &cli.BoolFlag{Name: "x", Usage: "extra functionality"}

func main() {
	script.Init(script.InitParams{
		Name:  "ezrun",
		Usage: script.ProcessUsageText(usageText),
	})
	app := &cli.App{
		Name:  "ez",
		Usage: "build and manage ezpkg.io project",
		Action: func(ctx *cli.Context) error {
			fmt.Println(`USAGE: ezrun COMMAND [ARGUMENT ...]
       ezrun --help`)
			return nil
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "verbose",
				Aliases:     []string{"v"},
				Usage:       "Turn on debug log",
				Destination: &flags.Verbose,
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "pkg",
				Usage:  "generate code for packages",
				Action: (&cmdPkg{}).Run,
				Flags:  []cli.Flag{flagAll},
			},
			{
				Name:   "list",
				Usage:  "list all packages",
				Action: (&cmdList{}).Run,
				Flags:  []cli.Flag{flagAll, flagX},
			},
			{
				Name:   "doc",
				Usage:  "generate README.md for all packages",
				Action: (&cmdDoc{}).Run,
				Flags:  []cli.Flag{flagAll},
			},
		},
	}
	initEnv()
	errorz.MustZ(app.Run(os.Args))
}

func initEnv() {
	env.EzpkgDir = mustDir(mustGetenv("EZPKG_DIR"))
	env.RootDir = mustDir(filepath.Dir(env.EzpkgDir))
	env.ZscriptsDir = mustDir(filepath.Join(env.EzpkgDir, "zscripts"))

	env.TargetDir = filepath.Join(env.RootDir, "ztarget")
	errorz.MustZ(os.MkdirAll(env.TargetDir, 0755))

	env.Info = mustReadJson[EzpkgInfo](filepath.Join(env.EzpkgDir, "ezpkg.json"))
	errorz.MustZ(env.Info.Validate())
}

func mustGetenv(key string) string {
	s := os.Getenv(key)
	if s == "" {
		panic(fmt.Sprintf("ENV %q is empty", key))
	}
	return s
}

func mustDir(path string) string {
	path = errorz.Must(filepath.Abs(path))
	dir := errorz.Must(os.Open(path))
	if errorz.Must(dir.Stat()).IsDir() {
		return path
	}
	panic(fmt.Sprintf("%q is not a directory", path))
}

func mustReadJson[T any](file string) (out *T) {
	data := errorz.Must(os.ReadFile(file))
	errorz.MustZ(json.Unmarshal(data, &out))
	return out
}
