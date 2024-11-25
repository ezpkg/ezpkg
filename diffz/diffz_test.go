package diffz_test

import (
	"fmt"
	"testing"

	. "ezpkg.io/conveyz"
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
	Convey("Diff", t, func() {
		Convey("ByChar", func() {
			Convey("default", func() {
				Convey("diff", func() {
					left, right := "onetwo threefour five", "onethree two fourfve"
					diffs := diffz.ByChar(left, right)
					Ω(diffs.IsDiff()).To(BeTrue())
				})
				Convey("multiline/diff", func() {
					diffs := diffz.ByChar(left0, right0)
					Ω(diffs.IsDiff()).To(BeTrue())
				})
			})
			Convey("ignore_space", func() {
				Convey("equal", func() {
					left, right := "onetwo threefour five", "onetwothree fourfive"
					diffs := diffz.IgnoreSpace().DiffByChar(left, right)
					Ω(diffs.IsDiff()).To(BeFalse())
				})
				Convey("diff", func() {
					left, right := "onetwo thre3four five", "onetwothree fourfive"
					diffs := diffz.IgnoreSpace().DiffByChar(left, right)
					Ω(diffs.IsDiff()).To(BeTrue())
				})
				Convey("multiline/equal", func() {
					diffs := diffz.IgnoreSpace().DiffByChar(left0, right0_IgnoreSpace)
					Ω(diffs.IsDiff()).To(BeFalse())
				})
			})
			Convey("placeholder", func() {
				Convey("equal", func() {
					left, right := "on█twothr██fourfiv█", "onetwothreefourfive"
					diffs := diffz.Placeholder().DiffByChar(left, right)
					Ω(diffs.IsDiff()).To(BeFalse())
				})
				Convey("space+placeholder/equal", func() {
					left, right := "on█two thr██four fiv█", "onetwothree fourfive"
					diffs := diffz.Placeholder().AndIgnoreSpace().DiffByChar(left, right)
					Ω(diffs.IsDiff()).To(BeFalse())
				})
				Convey("extra_x/equal", func() {
					left, right := "on█two thr██four fiv█", "onetwothreexfourfive"
					diffs := diffz.IgnoreSpace().AndPlaceholder().DiffByChar(left, right)
					Ω(diffs.IsDiff()).To(BeTrue())
				})
				Convey("extra_space/diff", func() {
					left, right := "on█two thr██four fiv█", "onetwothree fourfive"
					diffs := diffz.Placeholder().DiffByChar(left, right)
					Ω(diffs.IsDiff()).To(BeTrue())
				})
				Convey("extra_placeholder/diff", func() {
					left, right := "on█two thr███four fiv█", "onetwothree fourfive"
					diffs := diffz.Placeholder().AndIgnoreSpace().DiffByChar(left, right)
					Ω(diffs.IsDiff()).To(BeTrue())
				})
				Convey("multiline/equal", func() {
					diffs := diffz.IgnoreSpace().AndPlaceholder().DiffByChar(left0, right0_Placeholder)
					Ω(diffs.IsDiff()).To(BeFalse())
				})
				Convey("multiline/no_placeholder/diff", func() {
					diffs := diffz.IgnoreSpace().DiffByChar(left0, right0_Placeholder)
					Ω(diffs.IsDiff()).To(BeTrue())
				})
			})
		})
		Convey("ByLine", func() {
			Convey("diff/1", func() {
				diffs := diffz.ByLine("", "one")
				Ω(diffs.IsDiff()).To(BeTrue())
			})
			Convey("diff/2a", func() {
				diffs := diffz.ByLine("onetwo", "xonetwo")
				Ω(diffs.IsDiff()).To(BeTrue())
			})
			Convey("diff/2b", func() {
				diffs := diffz.ByLine("onetwo", "onextwo")
				Ω(diffs.IsDiff()).To(BeTrue())
			})
			Convey("diff/2c", func() {
				diffs := diffz.ByLine("onetwo", "onetwox")
				Ω(diffs.IsDiff()).To(BeTrue())
			})
			Convey("diff/3a", func() {
				diffs := diffz.ByLine("xonetwo", "onetwo")
				Ω(diffs.IsDiff()).To(BeTrue())
			})
			Convey("diff/3b", func() {
				diffs := diffz.ByLine("onextwo", "onetwo")
				Ω(diffs.IsDiff()).To(BeTrue())
			})
			Convey("diff/3c", func() {
				diffs := diffz.ByLine("onetwox", "onetwo")
				Ω(diffs.IsDiff()).To(BeTrue())
			})
			Convey("diff/4", func() {
				left, right := "zero", "one\ntwo\nthree"
				diffs := diffz.ByLine(left, right)
				Ω(diffs.IsDiff()).To(BeTrue())
				Ω(diffs.Items[0].Text).To(Equal("zero\n"), "❌0")
				Ω(diffs.Items[1].Text).To(Equal("one\n"), "❌1")
				Ω(diffs.Items[2].Text).To(Equal("two\n"), "❌2")
				Ω(diffs.Items[3].Text).To(Equal("three\n"), "❌3")
			})
			Convey("default", func() {
				diffs := diffz.ByLine(left0, right0)
				Ω(diffs.IsDiff()).To(BeTrue())
			})
			Convey("ignore_space", func() {
				Convey("not_ignore/diff", func() {
					diffs := diffz.ByLine(left0, right0_IgnoreSpace)
					Ω(diffs.IsDiff()).To(BeTrue())
				})
				Convey("ignore/equal", func() {
					diffs := diffz.IgnoreSpace().DiffByLine(left0, right0_IgnoreSpace)
					Ω(diffs.IsDiff()).To(BeFalse())
				})
				Convey("end_with_newline", func() {
					left, right := "zero", "one\ntwo\nthree"
					diffs := diffz.IgnoreSpace().DiffByLine(left, right)
					Ω(diffs.IsDiff()).To(BeTrue())
					Ω(diffs.Items[0].Text).To(Equal("zero\n"), "❌0")
					Ω(diffs.Items[1].Text).To(Equal("one\n"), "❌1")
					Ω(diffs.Items[2].Text).To(Equal("two\n"), "❌2")
					Ω(diffs.Items[3].Text).To(Equal("three\n"), "❌3")
				})
				Convey("color", func() {
					diffs := diffz.IgnoreSpace().DiffByLine(red, green)
					Ω(diffs.IsDiff()).To(BeTrue())
					Ω(diffs.Items[0].Type).To(Equal(diffz.DiffEqual), "❌0")
					Ω(diffs.Items[1].Type).To(Equal(diffz.DiffEqual), "❌1a")
					Ω(diffs.Items[1].Text).To(Equal("color:\n"), "❌1b")
					Ω(diffs.Items[2].Type).To(Equal(diffz.DiffDelete), "❌2a")
					Ω(diffs.Items[2].Text).To(Equal("  id: d56d5f0d-f05d-4d46-9ce2-af6396d25c55\n"), "❌2b")
					Ω(diffs.Items[3].Type).To(Equal(diffz.DiffInsert), "❌3a")
					Ω(diffs.Items[3].Text).To(Equal("  id: 5b01ec0b-0607-446e-8a25-aaef595902a9\n"), "❌3b")
					Ω(diffs.Items[4].Type).To(Equal(diffz.DiffDelete), "❌4a")
					Ω(diffs.Items[4].Text).To(Equal("  name: red\n"), "❌4b")
					Ω(diffs.Items[5].Type).To(Equal(diffz.DiffInsert), "❌5a")
					Ω(diffs.Items[5].Text).To(Equal("  name: green\n"), "❌5b")
					Ω(diffs.Items[6].Type).To(Equal(diffz.DiffEqual), "❌6a")
					Ω(diffs.Items[6].Text).To(Equal("  size: small\n"), "❌6b")
					Ω(diffs.Items[7].Type).To(Equal(diffz.DiffDelete), "❌7a")
					Ω(diffs.Items[7].Text).To(Equal("  code: #ff0000\n"), "❌7b")
					Ω(diffs.Items[8].Type).To(Equal(diffz.DiffInsert), "❌8a")
					Ω(diffs.Items[8].Text).To(Equal("  code: #00ff00\n"), "❌8b")
					Ω(len(diffs.Items)).To(Equal(9), "❌len")
				})
			})
			Convey("placeholder", func() {
				Convey("no_placeholder/diff", func() {
					diffs := diffz.ByLine(left0, right0_Placeholder)
					Ω(diffs.IsDiff()).To(BeTrue())
				})
				Convey("use_placeholder/equal", func() {
					diffs := diffz.Placeholder().AndIgnoreSpace().DiffByLine(left0, right0_Placeholder)
					Ω(diffs.IsDiff()).To(BeFalse())
				})
				Convey("color/equal", func() {
					diffs := diffz.IgnoreSpace().AndPlaceholder().DiffByLine(red, colorExpect)
					Ω(diffs.IsDiff()).To(BeFalse())
					Ω(diffs.Items[0].Text).To(Equal("\n"), "❌0")
					Ω(diffs.Items[1].Text).To(Equal("color:\n"), "❌1")
					Ω(diffs.Items[2].Text).To(Equal("  id: d56d5f0d-f05d-4d46-9ce2-af6396d25c55\n"), "❌2")
					Ω(diffs.Items[3].Text).To(Equal("  name: red\n"), "❌3")
					Ω(diffs.Items[4].Text).To(Equal("  size: small\n"), "❌4")
					Ω(diffs.Items[5].Text).To(Equal("  code: #ff0000\n"), "❌5")
				})
			})
		})
		Convey("regexp", func() {
			byChar := diffz.Placeholder().DiffByChar
			byLine := diffz.Placeholder().DiffByLine
			Convey("equal", func() {
				diffs := byChar("one42three", `one【\d+】three`)
				fmt.Println(diffz.Format(diffs))
				Ω(diffs.IsDiff()).To(BeFalse())
			})
			Convey("equal (trim space)", func() {
				diffs := byChar("one42three", `one【 \d+ 】three`)
				fmt.Println(diffz.Format(diffs))
				Ω(diffs.IsDiff()).To(BeFalse())
			})
			Convey("equal (match all, left)", func() {
				Convey("byChar", func() {
					diffs := byChar(`one【.+】three`, "one42three")
					fmt.Println(diffz.Format(diffs))
					Ω(diffs.IsDiff()).To(BeFalse())
				})
				Convey("byLine", func() {
					diffs := byLine("one\n【.+】three\nz", "one\n42three\nz")
					fmt.Println(diffz.Format(diffs))
					Ω(diffs.IsDiff()).To(BeFalse())
				})
				Convey("byLine (two)", func() {
					diffs := byLine("one\ntwo【.+】three\nz", "one\ntwo42three\nz")
					fmt.Println(diffz.Format(diffs))
					Ω(diffs.IsDiff()).To(BeFalse())
				})
				Convey("byLine (dot)", func() {
					diffs := byLine("one\ntwo【.+】three\nz", "one\ntwo.1.two42three\nz")
					fmt.Println(diffz.Format(diffs))
					Ω(diffs.IsDiff()).To(BeFalse())
				})
			})
			Convey("equal (match all, right)", func() {
				Convey("byChar", func() {
					diffs := byChar("one42three", `one【.+】three`)
					fmt.Println(diffz.Format(diffs))
					Ω(diffs.IsDiff()).To(BeFalse())
				})
				Convey("byLine", func() {
					diffs := byLine("one\n42three\nz", "one\n【 .+ 】three\nz")
					fmt.Println(diffz.Format(diffs))
					Ω(diffs.IsDiff()).To(BeFalse())
				})
				Convey("byLine (two)", func() {
					diffs := byLine("one\ntwo42three\nz", "one\ntwo【 .+ 】three\nz")
					fmt.Println(diffz.Format(diffs))
					Ω(diffs.IsDiff()).To(BeFalse())
				})
				Convey("byLine (dot)", func() {
					diffs := byLine("one\ntwo.1.two42three\nz", "one\ntwo【.+】three\nz")
					fmt.Println(diffz.Format(diffs))
					Ω(diffs.IsDiff()).To(BeFalse())
				})
			})
			Convey("equal (match all, both regexp)", func() {
				Convey("byChar", func() {
					diffs := byChar(`one【.+】three`, `one【.+】three`)
					fmt.Println(diffz.Format(diffs))
					Ω(diffs.IsDiff()).To(BeFalse())
				})
				Convey("byLine", func() {
					diffs := byLine("one\n【.+】three\nz", "one\n【.+】three\nz")
					fmt.Println(diffz.Format(diffs))
					Ω(diffs.IsDiff()).To(BeFalse())
				})
			})
			Convey("diff", func() {
				diffs := byChar("one4.2three", `one【[0-9]+】three`)
				fmt.Println(diffz.Format(diffs))
				Ω(diffs.IsDiff()).To(BeTrue())
			})
			Convey("diff (error)", func() {
				diffs := byChar("one4.2three", `one【***】three`)
				fmt.Println(diffz.Format(diffs))
				Ω(diffs.IsDiff()).To(BeTrue())
				Ω(diffz.Format(diffs)).To(ContainSubstring("❌《REGEXP:INVALID》"))
			})
		})
	})
}
