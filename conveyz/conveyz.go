// Package conveyz extends the package [convey] with additional functionality and make it work with [gomega].
//
// [convey]: github.com/smartystreets/goconvey/convey
// [gomega]: github.com/onsi/gomega
package conveyz // import "ezpkg.io/conveyz"

import (
	"strings"
	"testing"

	"github.com/onsi/gomega"
	gomegatypes "github.com/onsi/gomega/types"
	"github.com/smartystreets/goconvey/convey"

	"ezpkg.io/colorz"
	"ezpkg.io/fmtz"
	"ezpkg.io/stacktracez"
	"ezpkg.io/stringz"
)

func Convey(items ...any) {
	setupGomega(items...)
	convey.Convey(items...)
}
func SConvey(items ...any) {
	convey.SkipConvey(items...)
}
func SkipConvey(items ...any) {
	convey.SkipConvey(items...)
}
func FConvey(items ...any) {
	convey.FocusConvey(items...)
}
func FocusConvey(items ...any) {
	convey.FocusConvey(items...)
}
func Reset(action func()) {
	convey.Reset(action)
}

func GomegaExpect(actual any, extra ...any) gomega.Assertion {
	assertion := gomega.Expect(actual, extra...)
	return gomegaAssertion{actual: actual, assertion: assertion}
}

type gomegaAssertion struct {
	actual    any
	offset    int
	assertion gomega.Assertion
}

func (a gomegaAssertion) Should(matcher gomegatypes.GomegaMatcher, optionalDescription ...any) bool {
	return a.To(matcher, optionalDescription...)
}

func (a gomegaAssertion) ShouldNot(matcher gomegatypes.GomegaMatcher, optionalDescription ...any) bool {
	return a.ToNot(matcher, optionalDescription...)
}

func (a gomegaAssertion) To(matcher gomegatypes.GomegaMatcher, optionalDescription ...any) bool {
	convey.So(a.actual, func(_ any, _ ...any) string {
		success, err := matcher.Match(a.actual)
		if err != nil {
			stack := stacktracez.StackTraceSkip(4)
			return formatMsg(optionalDescription, stack, "%vUNEXPECTED: %v%v", colorz.Red, err, colorz.Yellow)
		}
		if !success {
			stack := stacktracez.StackTraceSkip(4)
			msg := matcher.FailureMessage(a.actual)
			return formatMsg(optionalDescription, stack, "%s", msg)
		}
		return ""
	})
	return true
}

func (a gomegaAssertion) ToNot(matcher gomegatypes.GomegaMatcher, optionalDescription ...any) bool {
	convey.So(a.actual, func(_ any, _ ...any) string {
		success, err := matcher.Match(a.actual)
		if err != nil {
			stack := stacktracez.StackTraceSkip(4)
			return formatMsg(optionalDescription, stack, "%vUNEXPECTED: %v%v", colorz.Red, err, colorz.Yellow)
		}
		if success {
			stack := stacktracez.StackTraceSkip(4)
			msg := matcher.NegatedFailureMessage(a.actual)
			return formatMsg(optionalDescription, stack, "%s", msg)
		}
		return ""
	})
	return true
}

func (a gomegaAssertion) NotTo(matcher gomegatypes.GomegaMatcher, optionalDescription ...any) bool {
	return a.ToNot(matcher, optionalDescription...)
}

func (a gomegaAssertion) WithOffset(offset int) gomegatypes.Assertion {
	return gomegaAssertion{
		actual:    a.actual,
		offset:    a.offset + offset,
		assertion: a.assertion,
	}
}

func (a gomegaAssertion) Error() gomegatypes.Assertion {
	return gomegaAssertion{
		actual:    a.actual,
		offset:    a.offset,
		assertion: a.assertion.Error(),
	}
}

func setupGomega(items ...any) {
	if len(items) >= 2 {
		testT, ok := items[1].(*testing.T)
		if ok {
			gomega.Default = gomega.NewWithT(testT)
		}
	}
}

func formatMsg(optionalDescription []any, stack *stacktracez.Frames, format string, args ...any) string {
	b := &stringz.Builder{}
	if len(optionalDescription) > 0 {
		b.Println(fmtz.FormatMsgArgsX(optionalDescription))
	}
	b.Printf(format, args...)
	b.Printf("\n\n")
	for _, frame := range stack.GetFrames() {
		pkg, _, _, _ := frame.Components()
		if pkg == "ezpkg.io/conveyz" ||
			strings.HasPrefix(pkg, "github.com/jtolds/gls") ||
			strings.HasPrefix(pkg, "github.com/smartystreets/goconvey") {
			continue
		}
		b.Printf("%s\n", frame)
	}
	return b.String()
}
