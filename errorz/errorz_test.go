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
			str := fmt.Sprintf(zErr.Error())
			assertEqual(t, str, `bar/one: foo`)
		})
		t.Run("print", func(t *testing.T) {
			str := fmt.Sprintln(zErr)
			assertEqual(t, str, `bar/one: foo`)
		})
		t.Run("printf:%s", func(t *testing.T) {
			str := fmt.Sprintf("%s\n", zErr)
			assertEqual(t, str, `bar/one: foo`)
		})
		t.Run("printf:%+v", func(t *testing.T) {
			str := fmt.Sprintf("%+v\n", zErr)
			assertEqual(t, str, `bar/one: foo`)
		})
		t.Run("printf:%v", func(t *testing.T) {
			str := fmt.Sprintf("%v\n", zErr)
			assertEqual(t, str, `bar/one: foo`)
		})
		t.Run("printf:%q", func(t *testing.T) {
			str := fmt.Sprintf("%q\n", zErr)
			assertEqual(t, str, `bar/one: foo`)
		})
	})
	t.Run("NoStack", func(t *testing.T) {
		err := errorz.Newf("foo/%v", "zero")
		zErr := errorz.NoStack().Wrapf(err, "bar/%v", "one")
		t.Run("printf:%s", func(t *testing.T) {
			str := fmt.Sprintf("%s\n", zErr)
			assertEqual(t, str, `bar/one: foo`)
		})
		t.Run("printf:%+v", func(t *testing.T) {
			str := fmt.Sprintf("%+v\n", zErr)
			assertEqual(t, str, `bar/one: foo`)
		})
	})
}

func assertEqual(t *testing.T, actual, expect string) {
	diffs := diffz.ByLineZ(actual, expect)
	if diffs.IsDiff() {
		fmt.Println(actual)
		fmt.Println()
		fmt.Println(diffz.Format(diffs))
		t.Error("error")
	} else {
		fmt.Println(actual)
	}
}
