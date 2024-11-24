package slicez_test

import (
	"fmt"
	"slices"
	"testing"

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
		Convey("Concat", func() {
			Ω(Concat([]int{1, 2}, []int{3, 4})).To(Equal([]int{1, 2, 3, 4}))
		})
		Convey("AppendTo", func() {
			s := []int{1, 2}
			sx := AppendTo(&s, 3, 4)
			Ω(s).To(Equal([]int{1, 2, 3, 4}))
			Ω(sx).To(Equal(s))
		})
		Convey("PrependTo", func() {
			s := []int{3, 4}
			sx := PrependTo(&s, 1, 2)
			Ω(s).To(Equal([]int{1, 2, 3, 4}))
			Ω(sx).To(Equal(s))
		})
		Convey("Unique", func() {
			Convey("empty", func() {
				s := []int{}
				fn := func(int) int { return 0 }

				Ω(IsUnique(s)).To(BeTrue())
				Ω(IsUniqueFunc(s, fn)).To(BeTrue())

				Ω(Unique(s)).To(Equal(s))
				Ω(UniqueFunc(s, fn)).To(Equal(s))
			})
			Convey("unique", func() {
				s := []int{1, 2, 3}
				fn := func(x int) int { return x }

				Ω(IsUnique(s)).To(BeTrue())
				Ω(IsUniqueFunc(s, fn)).To(BeTrue())

				Ω(Unique(s)).To(Equal(s))
				Ω(UniqueFunc(s, fn)).To(Equal(s))
			})
			Convey("not unique", func() {
				ss := [][]int{
					{3, 2, 1, 2},
					{2, 1, 3, 3},
					{1, 3, 1, 2, 1},
				}
				exps := [][]int{
					{3, 2, 1},
					{2, 1, 3},
					{1, 3, 2},
				}
				fn := func(x int) int { return x }
				for i, s := range ss {
					Convey(fmt.Sprintf("%d", i), func() {
						Ω(IsUnique(s)).To(BeFalse())
						Ω(IsUniqueFunc(s, fn)).To(BeFalse())

						Ω(Unique(s)).To(Equal(exps[i]))
						Ω(UniqueFunc(s, fn)).To(Equal(exps[i]))
					})
				}
			})
			Convey("sorted", func() {
				ss := [][]int{
					{3, 2, 1, 2},
					{2, 1, 3, 3},
					{1, 3, 1, 2, 1},
				}
				fn := func(x int) int { return x }
				for i, s := range ss {
					Convey(fmt.Sprintf("%d", i), func() {
						Ω(SortUnique(slices.Clone(s))).To(Equal([]int{1, 2, 3}))
						Ω(SortUniqueFunc(slices.Clone(s), fn)).To(Equal([]int{1, 2, 3}))
					})
				}
			})
		})
	})
}
