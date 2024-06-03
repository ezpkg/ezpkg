package stacktracez

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func TestStackTrace(t *testing.T) {
	formatStack := func(format string) string {
		st := StackTrace()
		return fmt.Sprintf(format, st)
	}
	st0, st1 := formatStack("%+v"), formatStack("%v")
	fmt.Println(st0)
	fmt.Println(st1)

	assert(t, regexp.MustCompile(`^github\.com/ezpkg/stacktracez\.StackTrace\n`).MatchString(st0)).
		Errorf("malform stacktrace")
	assert(t, regexp.MustCompile(`\n\t[\w/]+/ezpkg/ezpkg/stacktracez/stacktracez.go:\d+\n`).MatchString(st0)).
		Errorf("malform stacktrace")

	expected := regexp.MustCompile(strings.TrimSpace(`
github\.com/ezpkg/stacktracez/stacktracez\.go:\d+ · StackTrace
github\.com/ezpkg/stacktracez/stacktracez_test\.go:\d+ · TestStackTrace\.func1
github\.com/ezpkg/stacktracez/stacktracez_test\.go:\d+ · TestStackTrace
testing/testing\.go:\d+ · tRunner
`))
	assert(t, expected.MatchString(st1)).Errorf("malform stacktrace")
}

type assertFn func(format string, args ...any)

func (fn assertFn) Errorf(format string, args ...any) { fn(format, args...) }

func assert(t *testing.T, cond bool) assertFn {
	if cond {
		return func(string, ...any) {}
	} else {
		return t.Errorf
	}
}
