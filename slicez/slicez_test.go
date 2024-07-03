package slicez_test

import (
	"testing"

	. "github.com/onsi/gomega"

	. "ezpkg.io/conveyz"
	. "ezpkg.io/slicez"
)

func Test(t *testing.T) {
	Ω := GomegaExpect
	Convey("slicez", t, func() {
		slice := []string{"A", "B", "C", "D"}

		Convey("get", func() {
			Ω(GetX(slice, 0)).To(Equal("A"))
			Ω(GetX(slice, 2)).To(Equal("C"))
			Ω(GetX(slice, 5)).To(Equal(""))
			Ω(GetX(slice, -1)).To(Equal("D"))
			Ω(GetX(slice, -3)).To(Equal("B"))
			Ω(GetX(slice, -5)).To(Equal(""))
		})
		Convey("first", func() {
			Ω(First([]int{})).To(Equal(0))
			Ω(First(slice)).To(Equal("A"))

			Ω(FirstN([]int{}, 2)).To(Equal([]int{}))
			Ω(FirstN(slice, 0)).To(Equal([]string{}))
			Ω(FirstN(slice, 1)).To(Equal([]string{"A"}))
			Ω(FirstN(slice, -1)).To(Equal([]string{"D"}))
			Ω(FirstN(slice, 3)).To(Equal([]string{"A", "B", "C"}))
			Ω(FirstN(slice, -3)).To(Equal([]string{"B", "C", "D"}))
			Ω(FirstN(slice, 5)).To(Equal([]string{"A", "B", "C", "D"}))
			Ω(FirstN(slice, -5)).To(Equal([]string{"A", "B", "C", "D"}))
		})
		Convey("last", func() {
			Ω(Last(slice)).To(Equal("D"))

			Ω(LastN(slice, 1)).To(Equal([]string{"D"}))
			Ω(LastN(slice, -1)).To(Equal([]string{"A"}))
			Ω(LastN(slice, 5)).To(Equal([]string{"A", "B", "C", "D"}))
			Ω(LastN(slice, -5)).To(Equal([]string{"A", "B", "C", "D"}))
		})
		Convey("..Func", func() {
			fn := func(s string) bool {
				return s == "B" || s == "C"
			}
			Convey("FirstFunc", func() {
				x := FirstFunc(slice, fn)
				Ω(x).To(Equal("B"))
			})
			Convey("LastFunc", func() {
				x := LastFunc(slice, fn)
				Ω(x).To(Equal("C"))
			})
		})
	})
}
