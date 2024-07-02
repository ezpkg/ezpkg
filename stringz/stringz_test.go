package stringz_test

import (
	"testing"

	. "github.com/onsi/gomega"

	. "ezpkg.io/conveyz"
	"ezpkg.io/stringz"
)

func Test(t *testing.T) {
	Ω := GomegaExpect
	Convey("stringz", t, func() {
		input := "one,two,three,four"
		Convey("FirstNParts", func() {
			Convey("0", func() {
				x := stringz.FirstNParts(input, 0, ',')
				Ω(x).To(Equal(""))
			})
			Convey("2", func() {
				x := stringz.FirstNParts(input, 2, ',')
				Ω(x).To(Equal("one,two"))
			})
			Convey("4", func() {
				x := stringz.FirstNParts(input, 4, ',')
				Ω(x).To(Equal("one,two,three,four"))
			})
			Convey("10", func() {
				x := stringz.FirstNParts(input, 10, ',')
				Ω(x).To(Equal("one,two,three,four"))
			})
		})
		Convey("LastNParts", func() {
			Convey("0", func() {
				x := stringz.LastNParts(input, 0, ',')
				Ω(x).To(Equal(""))
			})
			Convey("3", func() {
				x := stringz.LastNParts(input, 3, ',')
				Ω(x).To(Equal("two,three,four"))
			})
		})
	})
}
