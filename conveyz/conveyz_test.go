package conveyz_test

import (
	"fmt"
	"testing"

	"github.com/onsi/gomega"

	"ezpkg.io/conveyz"
)

func Test(t *testing.T) {
	Ω := conveyz.Expect
	conveyz.Convey("Start", t, func() {
		s := "[0]"
		defer func() { fmt.Printf("\n%s\n", s) }()

		conveyz.Convey("Test 1:", func() {
			s += " → [1]"
			Ω(s).To(gomega.Equal("[0] → [1]"))

			conveyz.Convey("Test 1.1:", func() {
				s += " → [1.1]"
				Ω(s).To(gomega.Equal("[0] → [1] → [1.1]"))
			})
			conveyz.Convey("Test 1.2:", func() {
				s += " → [1.2]"
				Ω(s).To(gomega.Equal("[0] → [1] → [1.2]"))
			})
		})
		conveyz.Convey("Test 2:", func() {
			s += " → [2]"
			Ω(s).To(gomega.Equal("[0] → [2]"))

			conveyz.Convey("Test 2.1:", func() {
				s += " → [2.1]"
				Ω(s).To(gomega.Equal("[0] → [2] → [2.1]"))
			})
			conveyz.Convey("Test 2.2:", func() {
				s += " → [2.2]"
				Ω(s).To(gomega.Equal("[0] → [2] → [2.2]"))
			})
		})
		conveyz.SkipConvey("failure message", func() {
			// 👆 change SkipConvey to Convey to see failure messages

			conveyz.Convey("fail", func() {
				//  Expected
				//      <string>: [0] → [2]
				//  to equal
				//      <string>: this test will fail
				Ω(s).To(gomega.Equal("this test will fail"))
			})
			conveyz.Convey("UNEXPECTED ERROR", func() {
				// UNEXPECTED ERROR: Refusing to compare <nil> to <nil>.
				//  Be explicit and use BeNil() instead.  This is to avoid mistakes where both sides of an assertion are erroneously uninitialized.
				Ω(nil).To(gomega.Equal(nil))
			})
		})
	})
}
