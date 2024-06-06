package errorz_test

import (
	"fmt"
	"testing"

	"github.com/ezpkg/diffz"
	"github.com/ezpkg/errorz"
)

func TestError(t *testing.T) {
	t.Run("Wrap", func(t *testing.T) {
		err := fmt.Errorf("foo")
		zErr := errorz.Wrapf(err, "bar/%v", "one")
		t.Run("error", func(t *testing.T) {
			str := fmt.Sprint(zErr.Error())
			assertEqual(t, str, `bar/one: foo`)
		})
		t.Run("print", func(t *testing.T) {
			str := fmt.Sprint(zErr)
			assertEqual(t, str, `bar/one: foo`)
		})
		t.Run("printf:%s", func(t *testing.T) {
			str := fmt.Sprintf("%s", zErr)
			assertEqual(t, str, `bar/one: foo`)
		})
		t.Run("printf:%+v", func(t *testing.T) {
			str := fmt.Sprintf("%+v", zErr)
			fmt.Println(str)
			assertEqual(t, str, `
bar/one: foo
github.com/ezpkg/errorz_test.TestError.func1
	/Users/i/ws/ezpkg/ezpkg/errorz/errorz_test.go:██
testing.tRunner
	/usr/local/go/src/testing/testing.go:████
`)
		})
		t.Run("printf:%#v", func(t *testing.T) {
			str := fmt.Sprintf("%#v", zErr)
			fmt.Println(str)
			assertEqual(t, str, `
bar/one: foo
github.com/ezpkg/errorz_test/errorz_test.go:██ · TestError.func1
testing/testing.go:████ · tRunner
`)
		})
		t.Run("printf:%v", func(t *testing.T) {
			str := fmt.Sprintf("%v", zErr)
			assertEqual(t, str, `bar/one: foo`)
		})
		t.Run("printf:%q", func(t *testing.T) {
			str := fmt.Sprintf("%q", zErr)
			assertEqual(t, str, `"bar/one"`)
		})
	})
	t.Run("NoStack", func(t *testing.T) {
		err := errorz.Newf("foo/%v", "zero")
		zErr := errorz.NoStack().Wrapf(err, "bar/%v", "one")
		t.Run("printf:%s", func(t *testing.T) {
			str := fmt.Sprintf("%s", zErr)
			assertEqual(t, str, `bar/one: foo/zero`)
		})
		t.Run("printf:%+v", func(t *testing.T) {
			str := fmt.Sprintf("%+v", zErr)
			assertEqual(t, str, `bar/one: foo/zero`)
		})
	})
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

func assertEqual(t *testing.T, actual, expect string) {
	diffs := diffz.ByLineZ(actual, expect)
	if diffs.IsDiff() {
		fmt.Println(actual)
		fmt.Println()
		fmt.Println(diffz.Format(diffs))
		t.Error("❌ not equal")
	} else {
		fmt.Println(actual)
	}
}
