package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"

	"ezpkg.io/colorz"
	"ezpkg.io/errorz"
	"ezpkg.io/genz"
	"ezpkg.io/logz"
)

const flagVerbose = "verbose"

const helpText = `
EXAMPLES:
	  zgen codez-matchers ./...
`

func main() {
	cli.AppHelpTemplate = fmt.Sprintf("%s\n%s\n%s\n", colorz.Reset, cli.AppHelpTemplate, strings.TrimSpace(helpText))

	app := &cli.App{
		Name: "goz",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "verbose",
				Usage:   "enable verbosity (0: info, 4: debug, 8: more debug)",
				Value:   -int(logz.LevelWarn),
				Aliases: []string{"v"},
			},
		},
		Action: func(cx *cli.Context) error {
			cli.ShowAppHelpAndExit(cx, 0)
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:   "codez-matchers",
				Action: cmdCodezMatchers,
			},
		},
	}
	errorz.MustZ(app.Run(os.Args))
}

func initConfig(cx *cli.Context, plugins ...genz.Plugin) genz.Config {
	verbose := cx.Int(flagVerbose)
	if !cx.IsSet(flagVerbose) {
		env := os.Getenv(strings.ToUpper(flagVerbose))
		if env != "" {
			verbose, _ = strconv.Atoi(env)
		}
	}
	opt := &logz.TextHandlerOptions{
		AddSource:   false,
		Level:       -logz.Level(verbose),
		FormatLevel: logz.FormatLevelColor(nil),
	}
	logger := logz.New(logz.NewTextHandler(os.Stderr, opt))

	cfg := genz.Config{
		Logger:        logger,
		CleanOnly:     false,
		Namespace:     "",
		GoimportsArgs: []string{}, // example: -local github.com/foo
	}
	cfg.RegisterPlugin(plugins...)
	return cfg
}
