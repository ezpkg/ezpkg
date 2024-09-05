package testingz // import "ezpkg.io/testingz"

import (
	"fmt"
	"strings"

	"github.com/smartystreets/goconvey/convey"

	"ezpkg.io/colorz"
	"ezpkg.io/diffz"
	"ezpkg.io/fmtz"
	"ezpkg.io/typez"
)

func DiffByChar(expect, actual string) (formatted string, isDiff bool) {
	diffs := diffz.ByChar(actual, expect)
	return diffz.Format(diffs), diffs.IsDiff()
}

func DiffByCharX(expect, actual string, opt diffz.Option) (formatted string, isDiff bool) {
	diffs := diffz.ByCharX(actual, expect, opt)
	return diffz.Format(diffs), diffs.IsDiff()
}

func DiffByCharZ(expect, actual string) (formatted string, isDiff bool) {
	diffs := diffz.ByCharZ(actual, expect)
	return diffz.Format(diffs), diffs.IsDiff()
}

func DiffByLine(expect, actual string) (formatted string, isDiff bool) {
	diffs := diffz.ByLine(actual, expect)
	return diffz.Format(diffs), diffs.IsDiff()
}

func DiffByLineX(expect, actual string, opt diffz.Option) (formatted string, isDiff bool) {
	diffs := diffz.ByLineX(actual, expect, opt)
	return diffz.Format(diffs), diffs.IsDiff()
}

func DiffByLineZ(expect, actual string) (formatted string, isDiff bool) {
	diffs := diffz.ByLineZ(actual, expect)
	return diffz.Format(diffs), diffs.IsDiff()
}

// Usage with conveyz:
//
//	ΩxNoDiff := ConveyDiffByLine(diffz.IgnoreSpace().AndPlaceholder())
//	ΩxNoDiff(expect, actual, "my message")
func ConveyDiffByLine(opt diffz.Option) func(expect, actual string, msgArgs ...any) {
	pr := func(text string) {
		if opt.IgnoreSpace {
			fmt.Println(strings.TrimSpace(text))
		} else {
			fmt.Print(text)
			if !strings.HasSuffix(text, "\n") {
				fmt.Print(colorz.Yellow.Wrap("⛔\n(missing newline)\n"))
			}
		}
	}

	return func(expect, actual string, msgArgs ...any) {
		diffs := diffz.ByLineX(actual, expect, opt)
		if !diffs.IsDiff() {
			return
		}
		fmt.Print(colorz.Green.Wrap("\n👉 EXPECTED:\n"))
		pr(expect)
		fmt.Print(colorz.Red.Wrap("\n👉 ACTUAL:\n"))
		pr(actual)
		fmt.Print("\n👉 DIFF (", colorz.Red.Wrap("actual"), colorz.Green.Wrap("expected"), "):\n")
		fmt.Println(diffz.Format(diffs))

		msg := typez.Coalesce(fmtz.FormatMsgArgs(msgArgs), "unexpected diff")
		convey.So(0, func(any, ...any) string {
			return msg // failure with message
		})
	}
}

var _NoDiffByLine = ConveyDiffByLine(diffz.Option{})
var _NoDiffByLineZ = ConveyDiffByLine(diffz.IgnoreSpace().AndPlaceholder())

func ΩxNoDiffByLine(expect, actual string, msgArgs ...any) {
	_NoDiffByLine(expect, actual, msgArgs...)
}
func ΩxNoDiffByLineZ(expect, actual string, msgArgs ...any) {
	_NoDiffByLineZ(expect, actual, msgArgs...)
}
