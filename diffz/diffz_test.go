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

var right0_IgnoreSpace = `
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
}
`

var right0_Placeholder = `
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

func (██ Diffs) Unwrap() []Diff {

	return ██.Items
}
func (██ Diffs) IsDiff() bool {
	for _, █ := range ██.Items {

		if █.Type != DiffEqual {
			return true

		}
	}

	return true
}
`

func TestDiff(t *testing.T) {
	t.Run("ByChar", func(t *testing.T) {
		t.Run("default", func(t *testing.T) {
			left, right := "onetwo threefour five", "onethree two fourfve"
			diffs := diffz.ByChar(left, right)
			fmt.Println(diffz.Format(diffs))
			assert(t, diffs.IsDiff() == true).
				Errorf("expect diff")
		})
		t.Run("multiline", func(t *testing.T) {
			diffs := diffz.ByChar(left0, right0)
			fmt.Println(diffz.Format(diffs))
			assert(t, diffs.IsDiff() == true).
				Errorf("expect diff")
		})
		t.Run("ignore_space", func(t *testing.T) {
			t.Run("equal", func(t *testing.T) {
				left, right := "onetwo threefour five", "onetwothree fourfive"
				diffs := diffz.IgnoreSpace().DiffByChar(left, right)
				fmt.Println(diffz.Format(diffs))
			})
			t.Run("diff", func(t *testing.T) {
				left, right := "onetwo thre3four five", "onetwothree fourfive"
				diffs := diffz.IgnoreSpace().DiffByChar(left, right)
				fmt.Println(diffz.Format(diffs))
			})
			t.Run("multiline", func(t *testing.T) {
				diffs := diffz.IgnoreSpace().DiffByChar(left0, right0_IgnoreSpace)
				fmt.Println(diffz.Format(diffs))
				assert(t, diffs.IsDiff() == false).
					Errorf("expect no diff")
			})
		})
		t.Run("placeholder", func(t *testing.T) {
			t.Run("equal", func(t *testing.T) {
				left, right := "on█two thr██four fiv█", "onetwothreefourfive"
				diffs := diffz.Placeholder().DiffByChar(left, right)

				fmt.Println(diffz.Format(diffs))
				assert(t, diffs.IsDiff() == false).
					Errorf("expect no diff")
			})
			t.Run("equal/space+placeholder", func(t *testing.T) {
				left, right := "on█two thr██four fiv█", "onetwothree fourfive"
				diffs := diffz.Placeholder().AndIgnoreSpace().DiffByChar(left, right)

				fmt.Println(diffz.Format(diffs))
				assert(t, diffs.IsDiff() == false).
					Errorf("expect no diff")
			})
			t.Run("diff/extra_x", func(t *testing.T) {
				left, right := "on█two thr██four fiv█", "onetwothreexfourfive"
				diffs := diffz.IgnoreSpace().AndPlaceholder().DiffByChar(left, right)

				fmt.Println(diffz.Format(diffs))
				assert(t, diffs.IsDiff() == true).
					Errorf("expect diff")
			})
			t.Run("diff/extra_space", func(t *testing.T) {
				left, right := "on█two thr██four fiv█", "onetwothree fourfive"
				diffs := diffz.Placeholder().DiffByChar(left, right)

				fmt.Println(diffz.Format(diffs))
				assert(t, diffs.IsDiff() == true).
					Errorf("expect diff")
			})
			t.Run("diff/extra_placeholder", func(t *testing.T) {
				left, right := "on█two thr███four fiv█", "onetwothree fourfive"
				diffs := diffz.Placeholder().AndIgnoreSpace().DiffByChar(left, right)

				fmt.Println(diffz.Format(diffs))
				assert(t, diffs.IsDiff() == true).
					Errorf("expect diff")
			})
			t.Run("multiline/equal", func(t *testing.T) {
				diffs := diffz.IgnoreSpace().AndPlaceholder().DiffByChar(left0, right0_Placeholder)
				fmt.Println(diffz.Format(diffs))
				assert(t, diffs.IsDiff() == false).
					Errorf("expect no diff")
			})
			t.Run("multiline/no_placeholder", func(t *testing.T) {
				diffs := diffz.IgnoreSpace().DiffByChar(left0, right0_Placeholder)
				fmt.Println(diffz.Format(diffs))
				assert(t, diffs.IsDiff() == true).
					Errorf("expect diff")
			})
		})
	})
	t.Run("ByLine", func(t *testing.T) {
		t.Run("default", func(t *testing.T) {
			diffs := diffz.ByLine(left0, right0)
			fmt.Println(diffz.Format(diffs))
			assert(t, diffs.IsDiff() == true).
				Errorf("expect diff")
		})
		t.Run("ignore_space", func(t *testing.T) {
			t.Run("not_ignore/diff", func(t *testing.T) {
				diffs := diffz.ByLine(left0, right0_IgnoreSpace)
				fmt.Println(diffz.Format(diffs))
				assert(t, diffs.IsDiff() == true).
					Errorf("expect diff")
			})
			t.Run("ignore/equal", func(t *testing.T) {
				diffs := diffz.IgnoreSpace().DiffByLine(left0, right0_IgnoreSpace)
				fmt.Println(diffz.Format(diffs))
				assert(t, diffs.IsDiff() == false).
					Errorf("expect no diff")
			})
		})
		t.Run("placeholder", func(t *testing.T) {
			t.Run("no_placeholder/diff", func(t *testing.T) {
				diffs := diffz.ByLine(left0, right0_Placeholder)
				fmt.Println(diffz.Format(diffs))
				assert(t, diffs.IsDiff() == true).
					Errorf("expect diff")
			})
			t.Run("use_placeholder/equal", func(t *testing.T) {
				diffs := diffz.Placeholder().AndIgnoreSpace().DiffByLine(left0, right0_Placeholder)
				fmt.Println(diffz.Format(diffs))
				assert(t, diffs.IsDiff() == false).
					Errorf("expect no diff")
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
