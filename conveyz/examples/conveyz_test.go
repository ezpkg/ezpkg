package examples_test

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"

	"ezpkg.io/colorz"
	. "ezpkg.io/conveyz"
	"ezpkg.io/conveyz/examples"
)

func Test(t *testing.T) {
	Ω := GomegaExpect // 👈 make Ω as alias for GomegaExpect
	Convey("Start", t, func() {
		s := "[0]"
		defer func() { fmt.Printf("\n%s\n", s) }()

		add := func(part string) {
			s = examples.AppendStr(s, part)
		}

		Convey("Test 1:", func() {
			add(" → [1]")
			Ω(s).To(Equal("[0] → [1]"))

			Convey("Test 1.1:", func() {
				add(" → [1.1]")
				Ω(s).To(Equal("[0] → [1] → [1.1]"))
			})
			Convey("Test 1.2:", func() {
				add(" → [1.2]")
				Ω(s).To(Equal("[0] → [1] → [1.2]"))
			})
		})
		// 👇change to FConvey to focus on this block and all children
		// 👇change to SConvey to skip the block
		// 👇change to SConveyAsTODO to mark as TODO
		Convey("Test 2:", func() {
			add(" → [2]")
			Ω(s).To(Equal("[0] → [2]"))

			Convey("Test 2.1:", func() {
				add(" → [2.1]")
				Ω(s).To(Equal("[0] → [2] → [2.1]"))
			})
			Convey("Test 2.2:", func() {
				add(" → [2.2]")
				Ω(s).To(Equal("[0] → [2] → [2.2]"))
			})
		})
		SkipConveyAsTODO("failure message", func() {
			// 👆 change SkipConvey to Convey to see failure messages

			Convey(colorz.Cyan.Wrap("👉 this test will fail"), func() {
				//  Expected
				//      <string>: [0] → [2]
				//  to equal
				//      <string>: this test will fail
				Ω(s).To(Equal("this test will fail"))
			})
			Convey(colorz.Cyan.Wrap("👉 this test has UNEXPECTED error"), func() {
				// UNEXPECTED ERROR: Refusing to compare <nil> to <nil>.
				//  Be explicit and use BeNil() instead.  This is to avoid mistakes where both sides of an assertion are erroneously uninitialized.
				Ω(nil).To(Equal(nil))
			})
			Convey(colorz.Cyan.Wrap("👉 this test will panic"), func() {
				examples.CallFunc(func() {
					examples.WillPanic()
				})
			})
		})
	})
}
