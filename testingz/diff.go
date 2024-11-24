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

func DiffByChar(actual, expect string) (formatted string, isDiff bool) {
	diffs := diffz.ByChar(actual, expect)
	return diffz.Format(diffs), diffs.IsDiff()
}

func DiffByCharX(actual, expect string, opt diffz.Option) (formatted string, isDiff bool) {
	diffs := diffz.ByCharX(actual, expect, opt)
	return diffz.Format(diffs), diffs.IsDiff()
}

func DiffByCharZ(actual, expect string) (formatted string, isDiff bool) {
	diffs := diffz.ByCharZ(actual, expect)
	return diffz.Format(diffs), diffs.IsDiff()
}

func DiffByLine(actual, expect string) (formatted string, isDiff bool) {
	diffs := diffz.ByLine(actual, expect)
	return diffz.Format(diffs), diffs.IsDiff()
}

func DiffByLineX(actual, expect string, opt diffz.Option) (formatted string, isDiff bool) {
	diffs := diffz.ByLineX(actual, expect, opt)
	return diffz.Format(diffs), diffs.IsDiff()
}

func DiffByLineZ(actual, expect string) (formatted string, isDiff bool) {
	diffs := diffz.ByLineZ(actual, expect)
	return diffz.Format(diffs), diffs.IsDiff()
}

// Usage with conveyz:
//
//	ΩxNoDiff := ConveyDiffByLine(diffz.IgnoreSpace().AndPlaceholder())
//	ΩxNoDiff(expect, actual, "my message")
func ConveyDiffByLine(opt diffz.Option) func(actual, expect string, msgArgs ...any) {
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

	return func(actual, expect string, msgArgs ...any) {
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

func ΩxNoDiffByLine(actual, expect string, msgArgs ...any) {
	_NoDiffByLine(actual, expect, msgArgs...)
}
func ΩxNoDiffByLineZ(actual, expect string, msgArgs ...any) {
	_NoDiffByLineZ(actual, expect, msgArgs...)
}
