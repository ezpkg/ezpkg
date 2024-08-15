package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"ezpkg.io/errorz"
	"ezpkg.io/genz"
	"ezpkg.io/genz/plugins/sample"
	"ezpkg.io/logz"
	"ezpkg.io/typez"
)

var flClean = flag.Bool("clean", false, "clean generated files without generating new files")
var flPlugin = flag.String("plugin", "", "comma separated list of plugins for generating (default to all plugins)")
var flNamespace = flag.String("namespace", "", "github.com/myproject")
var flVerbose = flag.Int("v", -4, "enable verbosity (-4: warn, 0: info, 4: debug, 8: more debug)")

func usage() {
	const text = `
Usage: genz [OPTION] PATTERN ...

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

func Start(plugins ...genz.Plugin) {
	flag.Parse()
	patterns := flag.Args()
	if len(patterns) == 0 {
		usage()
		os.Exit(2)
	}

	opt := &logz.TextHandlerOptions{
		Level:       logz.Level(-typez.Deptr(flVerbose)),
		FormatLevel: logz.FormatLevelColor(nil),
	}
	logger := logz.New(logz.NewTextHandler(os.Stderr, opt))

	cfg := genz.Config{
		Logger:        logger,
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

	ctx := context.Background()
	errorz.MustZ(genz.Start(ctx, cfg, patterns...))
}
