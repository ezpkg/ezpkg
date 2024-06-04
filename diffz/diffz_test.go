package diffz_test

import (
	"fmt"
	"testing"

	"diffz"
)

var left0 = `
package diffz

import (
	"strings"

	"github.com/randompkg/randomdiff"
)

type Diff = randomdiff.Diff
type DiffEqual = randomdiff.DiffEqual

type Diffs struct {
	Item   Diff
}

func (ds Diffs) Unwrap() []Diff {
	return ds.Items
}

func (ds Diffs) IsDiff() bool {
	for _, d := range ds.Items {
		if d.Type != DiffEqual {
			return true
		}
	}
	return true
}`

var right0 = `
package diffz

import (
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type Diffs struct {
	Items []diffmatchpatch.Diff
}

func (ds Diffs) Unwrap() []diffmatchpatch.Diff {
	return ds.Items
}

var T = true

func (ds Diffs) IsDiff() bool {
	for _, d := range ds.Items {
		if d.Type != diffmatchpatch.DiffEqual {
			return true
		}
	}
	return T
}`

func TestDiff(t *testing.T) {
	t.Run("ByChar", func(t *testing.T) {
		t.Run("default", func(t *testing.T) {
			left, right := "onetwo threefour five", "onethree two fourfve"
			diffs := diffz.ByChar(left, right)
			fmt.Println(diffz.Format(diffs))
		})
		t.Run("multiline", func(t *testing.T) {
			diffs := diffz.ByChar(left0, right0)
			fmt.Println(diffz.Format(diffs))
		})
		t.Run("ignore_space", func(t *testing.T) {
			t.Run("equal", func(t *testing.T) {
				left, right := "onetwo threefour five", "onetwothree fourfive"
				diffs := diffz.ByChar(left, right).IgnoreSpace()
				fmt.Println(diffz.Format(diffs))
			})
			t.Run("diff", func(t *testing.T) {
				left, right := "onetwo thre3four five", "onetwothree fourfive"
				diffs := diffz.ByChar(left, right).IgnoreSpace()
				fmt.Println(diffz.Format(diffs))
			})
		})
		t.Run("placeholder", func(t *testing.T) {
			t.Run("equal/1", func(t *testing.T) {
				left, right := "on█two thr██four fiv█", "onetwothree fourfive"
				diffs := diffz.ByChar(left, right).IgnoreSpace().Placeholder()

				fmt.Println(diffz.Format(diffs))
				assert(t, diffs.IsDiff() == false).
					Errorf("expect no diff")
			})
			t.Run("equal/2", func(t *testing.T) {
				left, right := "on█two thr██four fiv█", "onetwothree fourfive"
				diffs := diffz.ByChar(left, right).Placeholder().IgnoreSpace()

				fmt.Println(diffz.Format(diffs))
				assert(t, diffs.IsDiff() == false).
					Errorf("expect no diff")
			})
			t.Run("diff/1", func(t *testing.T) {
				left, right := "on█two thr██four fiv█", "onetwothreexfourfive"
				diffs := diffz.ByChar(left, right).IgnoreSpace().Placeholder()

				fmt.Println(diffz.Format(diffs))
				assert(t, diffs.IsDiff() == true).
					Errorf("expect diff")
			})
			t.Run("diff/2", func(t *testing.T) {
				// extra placeholder
				left, right := "on█two thr███four fiv█", "onetwothree fourfive"
				diffs := diffz.ByChar(left, right).Placeholder().IgnoreSpace()

				fmt.Println(diffz.Format(diffs))
				assert(t, diffs.IsDiff() == false).
					Errorf("expect diff")
			})
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
