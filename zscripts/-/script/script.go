package script

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/urfave/cli/v2"

	"ezpkg.io/colorz"
)

var params InitOutput

type InitParams struct {
	Name  string
	Usage string
}

type InitOutput struct {
	InitParams
	Verbose bool
}

func (x *InitOutput) FlagVerbose() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:        "verbose",
		Aliases:     []string{"v"},
		Usage:       "Turn on debug log",
		Destination: &x.Verbose,
	}
}

func FlagVerbose() *cli.BoolFlag { return params.FlagVerbose() }

func Init(p InitParams) {
	cli.AppHelpTemplate = fmt.Sprintf("%s\n%s\n%s", colorz.Reset, cli.AppHelpTemplate, p.Usage)
	params = InitOutput{
		InitParams: p,
	}
}

type Args []string

func WrapArgs(cx *cli.Context) *Args {
	args := Args(cx.Args().Slice())
	return &args
}

func (args *Args) Len() int {
	return len(*args)
}

func (args *Args) IsEmpty() bool { return len(*args) == 0 }

func (args *Args) Slice() []string {
	return *args
}

func (args *Args) Get(i int) string {
	if 0 <= i && i < len(*args) {
		return (*args)[i]
	}
	return ""
}

func (args *Args) Next() string {
	return args.Get(0)
}

func (args *Args) Consume() string {
	if len(*args) == 0 {
		return ""
	}
	out := (*args)[0]
	*args = (*args)[1:]
	return out
}

func (args *Args) MustConsume(name string) string {
	if len(*args) == 0 {
		ExitWithUsagef("missing argument %v", name)
	}
	out := (*args)[0]
	*args = (*args)[1:]
	return out
}

func (args *Args) MustConsumeLast(name string) string {
	if len(*args) == 0 {
		ExitWithUsagef("missing argument %v", name)
	}
	out := (*args)[len(*args)-1]
	*args = (*args)[:len(*args)-1]
	return out
}

func (args *Args) ConsumeByFunc(fn func(arg string) bool) string {
	for i, arg := range *args {
		if fn(arg) {
			*args = append((*args)[:i], (*args)[i+1:]...)
			return arg
		}
	}
	return ""
}

func (args *Args) MustConsumeRemain(min int, name string) []string {
	if len(*args) < min {
		if min == 1 {
			ExitWithUsagef("missing argument %v", name)
		} else {
			ExitWithUsagef("missing arguments %v (expect %v)", name, min)
		}
	}
	out := args.Slice()
	*args = nil
	return out
}

func (args *Args) MustEmpty() {
	if args.Len() > 0 {
		ExitWithUsagef("unexpected argument %v", args.Get(0))
	}
}

func Exitf(msg string, args ...any) {
	msg = strings.ReplaceAll(msg, `\n`, "\n")
	Stderrf(msg, args...)
	Stderrf("\n\n")
	Stderrf(colorz.Reset.Code())
	os.Exit(1)
}

func ExitWithUsagef(msg string, args ...any) {
	if msg != "" {
		msg = strings.ReplaceAll(msg, `\n`, "\n")
		Stderrf(msg, args...)
		Stderrf("\n\n")
	}
	Stderrf(params.Usage)
	Stderrf("\n\n")
	Stderrf(colorz.Reset.Code())
	os.Exit(1)
}

func Stderrf(msg string, args ...any) {
	fprintf(os.Stderr, colorz.Yellow.Code())
	fprintf(os.Stderr, msg, args...)
	fprintf(os.Stderr, colorz.Reset.Code())
}

func fprintf(w io.Writer, format string, args ...any) {
	must(fmt.Fprintf(w, format, args...))
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
