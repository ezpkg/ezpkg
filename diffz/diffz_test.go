package diffz_test

import (
	"fmt"
	"testing"

	"ezpkg.io/diffz"
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

var colorExpect = `
color:
  id: ████████-████-████-████-████████████
  name: red
  size: small
  code: #ff0000`
var red = `
color:
  id: d56d5f0d-f05d-4d46-9ce2-af6396d25c55
  name: red
  size: small
  code: #ff0000`
var green = `
color:
  id: 5b01ec0b-0607-446e-8a25-aaef595902a9
  name: green
  size: small
  code: #00ff00`

func TestDiff(t *testing.T) {
	t.Run("ByChar", func(t *testing.T) {
		t.Run("default", func(t *testing.T) {
			t.Run("diff", func(t *testing.T) {
				left, right := "onetwo threefour five", "onethree two fourfve"
				diffs := diffz.ByChar(left, right)
				assertDiff(t, diffs)
			})
			t.Run("multiline/diff", func(t *testing.T) {
				diffs := diffz.ByChar(left0, right0)
				assertDiff(t, diffs)
			})
		})
		t.Run("ignore_space", func(t *testing.T) {
			t.Run("equal", func(t *testing.T) {
				left, right := "onetwo threefour five", "onetwothree fourfive"
				diffs := diffz.IgnoreSpace().DiffByChar(left, right)
				assertEqual(t, diffs)
			})
			t.Run("diff", func(t *testing.T) {
				left, right := "onetwo thre3four five", "onetwothree fourfive"
				diffs := diffz.IgnoreSpace().DiffByChar(left, right)
				assertDiff(t, diffs)
			})
			t.Run("multiline/equal", func(t *testing.T) {
				diffs := diffz.IgnoreSpace().DiffByChar(left0, right0_IgnoreSpace)
				assertEqual(t, diffs)
			})
		})
		t.Run("placeholder", func(t *testing.T) {
			t.Run("equal", func(t *testing.T) {
				left, right := "on█twothr██fourfiv█", "onetwothreefourfive"
				diffs := diffz.Placeholder().DiffByChar(left, right)
				assertEqual(t, diffs)
			})
			t.Run("space+placeholder/equal", func(t *testing.T) {
				left, right := "on█two thr██four fiv█", "onetwothree fourfive"
				diffs := diffz.Placeholder().AndIgnoreSpace().DiffByChar(left, right)
				assertEqual(t, diffs)
			})
			t.Run("extra_x/equal", func(t *testing.T) {
				left, right := "on█two thr██four fiv█", "onetwothreexfourfive"
				diffs := diffz.IgnoreSpace().AndPlaceholder().DiffByChar(left, right)
				assertDiff(t, diffs)
			})
			t.Run("extra_space/diff", func(t *testing.T) {
				left, right := "on█two thr██four fiv█", "onetwothree fourfive"
				diffs := diffz.Placeholder().DiffByChar(left, right)
				assertDiff(t, diffs)
			})
			t.Run("extra_placeholder/diff", func(t *testing.T) {
				left, right := "on█two thr███four fiv█", "onetwothree fourfive"
				diffs := diffz.Placeholder().AndIgnoreSpace().DiffByChar(left, right)
				assertDiff(t, diffs)
			})
			t.Run("multiline/equal", func(t *testing.T) {
				diffs := diffz.IgnoreSpace().AndPlaceholder().DiffByChar(left0, right0_Placeholder)
				assertEqual(t, diffs)
			})
			t.Run("multiline/no_placeholder/diff", func(t *testing.T) {
				diffs := diffz.IgnoreSpace().DiffByChar(left0, right0_Placeholder)
				assertDiff(t, diffs)
			})
		})
	})
	t.Run("ByLine", func(t *testing.T) {
		t.Run("diff/1", func(t *testing.T) {
			diffs := diffz.ByLine("", "one")
			assertDiff(t, diffs)
		})
		t.Run("diff/2a", func(t *testing.T) {
			diffs := diffz.ByLine("onetwo", "xonetwo")
			assertDiff(t, diffs)
		})
		t.Run("diff/2b", func(t *testing.T) {
			diffs := diffz.ByLine("onetwo", "onextwo")
			assertDiff(t, diffs)
		})
		t.Run("diff/2c", func(t *testing.T) {
			diffs := diffz.ByLine("onetwo", "onetwox")
			assertDiff(t, diffs)
		})
		t.Run("diff/3a", func(t *testing.T) {
			diffs := diffz.ByLine("xonetwo", "onetwo")
			assertDiff(t, diffs)
		})
		t.Run("diff/3b", func(t *testing.T) {
			diffs := diffz.ByLine("onextwo", "onetwo")
			assertDiff(t, diffs)
		})
		t.Run("diff/3c", func(t *testing.T) {
			diffs := diffz.ByLine("onetwox", "onetwo")
			assertDiff(t, diffs)
		})
		t.Run("diff/4", func(t *testing.T) {
			left, right := "zero", "one\ntwo\nthree"
			diffs := diffz.ByLine(left, right)
			assertDiff(t, diffs)
			assert(t, diffs.Items[0].Text == "zero\n").Errorf("❌0")
			assert(t, diffs.Items[1].Text == "one\n").Errorf("❌1")
			assert(t, diffs.Items[2].Text == "two\n").Errorf("❌2")
			assert(t, diffs.Items[3].Text == "three\n").Errorf("❌3")
		})
		t.Run("default", func(t *testing.T) {
			diffs := diffz.ByLine(left0, right0)
			assertDiff(t, diffs)
		})
		t.Run("ignore_space", func(t *testing.T) {
			t.Run("not_ignore/diff", func(t *testing.T) {
				diffs := diffz.ByLine(left0, right0_IgnoreSpace)
				assertDiff(t, diffs)
			})
			t.Run("ignore/equal", func(t *testing.T) {
				diffs := diffz.IgnoreSpace().DiffByLine(left0, right0_IgnoreSpace)
				assertEqual(t, diffs)
			})
			t.Run("end_with_newline", func(t *testing.T) {
				left, right := "zero", "one\ntwo\nthree"
				diffs := diffz.IgnoreSpace().DiffByLine(left, right)
				assertDiff(t, diffs)
				assert(t, diffs.Items[0].Text == "zero\n").Errorf("❌0")
				assert(t, diffs.Items[1].Text == "one\n").Errorf("❌1")
				assert(t, diffs.Items[2].Text == "two\n").Errorf("❌2")
				assert(t, diffs.Items[3].Text == "three\n").Errorf("❌3")
			})
			t.Run("color", func(t *testing.T) {
				diffs := diffz.IgnoreSpace().DiffByLine(red, green)
				assertDiff(t, diffs)
				assert(t, diffs.Items[0].Type == diffz.DiffEqual).Errorf("❌0")
				assert(t, diffs.Items[1].Type == diffz.DiffEqual).Errorf("❌1a")
				assert(t, diffs.Items[1].Text == "color:\n").Errorf("❌1b")
				assert(t, diffs.Items[2].Type == diffz.DiffDelete).Errorf("❌2a")
				assert(t, diffs.Items[2].Text == "  id: d56d5f0d-f05d-4d46-9ce2-af6396d25c55\n").Errorf("❌2b")
				assert(t, diffs.Items[3].Type == diffz.DiffInsert).Errorf("❌3a")
				assert(t, diffs.Items[3].Text == "  id: 5b01ec0b-0607-446e-8a25-aaef595902a9\n").Errorf("❌3b")
				assert(t, diffs.Items[4].Type == diffz.DiffDelete).Errorf("❌4a")
				assert(t, diffs.Items[4].Text == "  name: red\n").Errorf("❌4b")
				assert(t, diffs.Items[5].Type == diffz.DiffInsert).Errorf("❌5a")
				assert(t, diffs.Items[5].Text == "  name: green\n").Errorf("❌5b")
				assert(t, diffs.Items[6].Type == diffz.DiffEqual).Errorf("❌6a")
				assert(t, diffs.Items[6].Text == "  size: small\n").Errorf("❌6b")
				assert(t, diffs.Items[7].Type == diffz.DiffDelete).Errorf("❌7a")
				assert(t, diffs.Items[7].Text == "  code: #ff0000\n").Errorf("❌7b")
				assert(t, diffs.Items[8].Type == diffz.DiffInsert).Errorf("❌8a")
				assert(t, diffs.Items[8].Text == "  code: #00ff00\n").Errorf("❌8b")
				assert(t, len(diffs.Items) == 9).Errorf("❌len")
			})
		})
		t.Run("placeholder", func(t *testing.T) {
			t.Run("no_placeholder/diff", func(t *testing.T) {
				diffs := diffz.ByLine(left0, right0_Placeholder)
				assertDiff(t, diffs)
			})
			t.Run("use_placeholder/equal", func(t *testing.T) {
				diffs := diffz.Placeholder().AndIgnoreSpace().DiffByLine(left0, right0_Placeholder)
				assertEqual(t, diffs)
			})
			t.Run("color/equal", func(t *testing.T) {
				diffs := diffz.IgnoreSpace().AndPlaceholder().DiffByLine(red, colorExpect)
				assertEqual(t, diffs)
				assert(t, diffs.Items[0].Text == "\n").Errorf("❌0")
				assert(t, diffs.Items[1].Text == "color:\n").Errorf("❌1")
				assert(t, diffs.Items[2].Text == "  id: d56d5f0d-f05d-4d46-9ce2-af6396d25c55\n").Errorf("❌2")
				assert(t, diffs.Items[3].Text == "  name: red\n").Errorf("❌3")
				assert(t, diffs.Items[4].Text == "  size: small\n").Errorf("❌4")
				assert(t, diffs.Items[5].Text == "  code: #ff0000\n").Errorf("❌5")
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

func assertDiff(t *testing.T, diffs diffz.Diffs) {
	fmt.Println(diffz.Format(diffs))
	assert(t, diffs.IsDiff() == true).
		Errorf("expect diff")
}

func assertEqual(t *testing.T, diffs diffz.Diffs) {
	fmt.Println(diffz.Format(diffs))
	assert(t, diffs.IsDiff() == false).
		Errorf("expect no diff")
}
