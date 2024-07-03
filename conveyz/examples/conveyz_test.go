package examples_test

import (
	"fmt"
	"testing"

	"github.com/onsi/gomega"

	"ezpkg.io/colorz"
	"ezpkg.io/conveyz"
	"ezpkg.io/conveyz/examples"
)

func Test(t *testing.T) {
	Ω := conveyz.GomegaExpect
	conveyz.Convey("Start", t, func() {
		s := "[0]"
		defer func() { fmt.Printf("\n%s\n", s) }()

		add := func(part string) {
			s = examples.AppendStr(s, part)
		}

		conveyz.Convey("Test 1:", func() {
			add(" → [1]")
			Ω(s).To(gomega.Equal("[0] → [1]"))

			conveyz.Convey("Test 1.1:", func() {
				add(" → [1.1]")
				Ω(s).To(gomega.Equal("[0] → [1] → [1.1]"))
			})
			conveyz.Convey("Test 1.2:", func() {
				add(" → [1.2]")
				Ω(s).To(gomega.Equal("[0] → [1] → [1.2]"))
			})
		})
		conveyz.Convey("Test 2:", func() {
			add(" → [2]")
			Ω(s).To(gomega.Equal("[0] → [2]"))

			conveyz.Convey("Test 2.1:", func() {
				add(" → [2.1]")
				Ω(s).To(gomega.Equal("[0] → [2] → [2.1]"))
			})
			conveyz.Convey("Test 2.2:", func() {
				add(" → [2.2]")
				Ω(s).To(gomega.Equal("[0] → [2] → [2.2]"))
			})
		})
		conveyz.SkipConveyAsTODO("failure message", func() {
			// 👆 change SkipConvey to Convey to see failure messages

			conveyz.Convey(colorz.Cyan.Wrap("👉 this test will fail"), func() {
				//  Expected
				//      <string>: [0] → [2]
				//  to equal
				//      <string>: this test will fail
				Ω(s).To(gomega.Equal("this test will fail"))
			})
			conveyz.Convey(colorz.Cyan.Wrap("👉 this test has UNEXPECTED error"), func() {
				// UNEXPECTED ERROR: Refusing to compare <nil> to <nil>.
				//  Be explicit and use BeNil() instead.  This is to avoid mistakes where both sides of an assertion are erroneously uninitialized.
				Ω(nil).To(gomega.Equal(nil))
			})
			conveyz.Convey(colorz.Cyan.Wrap("👉 this test will panic"), func() {
				examples.CallFunc(func() {
					examples.WillPanic()
				})
			})
		})
	})
}