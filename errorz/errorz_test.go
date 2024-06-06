package errorz_test

import (
	"fmt"
	"regexp"
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
			assertRegexp(t, str, `^bar/one: foo\n`)
			assertRegexp(t, str, `github[.]com/ezpkg/errorz[.]Wrapf\n`)
			assertRegexp(t, str, `/ezpkg/ezpkg/errorz/errorz[.]go:\d+\n`)
			assertRegexp(t, str, `\ntesting[.]tRunner\n`)
			assertRegexp(t, str, `/testing/testing[.]go:\d+\n$`)
		})
		t.Run("printf:%#v", func(t *testing.T) {
			str := fmt.Sprintf("%#v", zErr)
			fmt.Println(str)
			assertRegexp(t, str, `^bar/one: foo\n`)
			assertRegexp(t, str, `\ngithub[.]com/ezpkg/errorz/errorz[.]go:\d+ · Wrapf\n`)
			assertRegexp(t, str, `\ntesting/testing[.]go:\d+ · tRunner\n$`)
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

func assertEqual(t *testing.T, actual, expect string) {
	diffs := diffz.ByLine(actual, expect)
	if diffs.IsDiff() {
		fmt.Println(actual)
		fmt.Println()
		fmt.Println(diffz.Format(diffs))
		t.Error("❌ not equal")
	} else {
		fmt.Println(actual)
	}
}

func assertRegexp(t *testing.T, actual, expect string) {
	if !regexp.MustCompile(expect).MatchString(actual) {
		t.Errorf("❌ expect %q to contain %q", actual, expect)
	}
}
