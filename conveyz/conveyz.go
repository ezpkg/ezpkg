// Package conveyz extends the package [convey] with additional functionality and make it work with [gomega].
//
// [convey]: github.com/smartystreets/goconvey/convey
// [gomega]: github.com/onsi/gomega
package conveyz // import "ezpkg.io/conveyz"

import (
	"fmt"
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

var skippedTests bool

func Convey(items ...any) {
	defer setupGomega(items...)()
	convey.Convey(items...)
}

// SConvey (alias of SkipConvey) skips the current scope and all child scopes. It also makes the test fail.
func SConvey(items ...any) {
	skippedTests=true
	patchMessage(items, "SKIP", colorz.Magenta)
	convey.SkipConvey(items...)
}

// SkipConvey skips the current scope and all child scopes. It also makes the test fail.
func SkipConvey(items ...any) {
	skippedTests=true
	patchMessage(items, "SKIP", colorz.Magenta)
	convey.SkipConvey(items...)
}

// FConvey (alias of FocusConvey) runs the current scope and all child scopes, but skips all other scopes. It also makes the test fail.
func FConvey(items ...any) {
	skippedTests=true
	patchMessage(items, "FOCUS", colorz.Magenta)
	convey.FocusConvey(items...)
}

// FocusConvey runs the current scope and all child scopes, but skips all other scopes. It also makes the test fail.
func FocusConvey(items ...any) {
	skippedTests=true
	patchMessage(items, "FOCUS", colorz.Magenta)
	convey.FocusConvey(items...)
}

// SkipConveyAsTODO is similar to SkipConvey but does not make the test fail.
func SkipConveyAsTODO(items ...any) {
	patchMessage(items, "TODO", colorz.Magenta)
	convey.SkipConvey(items...)
}

// Reset registers a cleanup function to run after each Convey() in the same scope.
func Reset(action func()) {
	convey.Reset(action)
}

// GomegaExpect is an adapter to make gomega work with goconvey.
//
// Usage: Ω := GomegaExpect
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

func setupGomega(items ...any) func() {
	if len(items) < 2 {
		return func() {}
	}
	testT, ok := items[1].(*testing.T)
	if !ok {
		return func() {}
	}
	// this is top-level convey, init gomega
	gomega.Default = gomega.NewWithT(testT)
	return func() {
		if skippedTests {
			fmt.Println(colorz.Magenta.Wrap("--- NOTE: There are skipped/focused tests. Make sure to include them or mark as TODO."))
			testT.Fail()
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

func patchMessageX(items []any, fn func(s string) string) {
	if len(items) == 0 {
		return
	}
	if msg, ok := items[0].(string); ok {
		items[0] = fn(msg)
	}
}

func patchMessage(items []any, prefix string, color colorz.Color) {
	patchMessageX(items, func(s string) string {
		return fmt.Sprintf("%s%s:%s %s", color, prefix, colorz.Reset, s)
	})
}
