package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"ezpkg.io/genz"
	"ezpkg.io/genz/plugins/sample"
)

var flClean = flag.Bool("clean", false, "clean generated files without generating new files")
var flPlugin = flag.String("plugin", "", "comma separated list of plugins for generating (default to all plugins)")
var flNamespace = flag.String("namespace", "", "github.com/myproject")
var flVerbose = flag.Int("verbose", 0, "enable verbosity (0: info, 4: debug, 8: more debug)")

func usage() {
	const text = `
Usage: ggen [OPTION] PATTERN ...

Options:
`
	fmt.Print(text[1:])
	flag.PrintDefaults()
}

func main() {
	Start(
		sample.New(), // sample plugin
	)
}

func Start(plugins ...ggen.Plugin) {
	flag.Parse()
	patterns := flag.Args()
	if len(patterns) == 0 {
		usage()
		os.Exit(2)
	}

	cfg := ggen.Config{
		LogLevel:      -ggen.LogLevel(*flVerbose),
		CleanOnly:     *flClean,
		Namespace:     *flNamespace,
		GoimportsArgs: []string{}, // example: -local github.com/foo
	}
	cfg.RegisterPlugin(plugins...)
	if *flPlugin != "" {
		pluginNames := strings.Split(*flPlugin, ",")
		for _, name := range pluginNames {
			cfg.EnablePlugin(name)
		}
	}

	must(ggen.Start(cfg, patterns...))
}

func must(err error) {
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}
