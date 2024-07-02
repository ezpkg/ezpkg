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
		Convey("empty", func() {
			Convey("FirstNParts", func() {
				x, remain := stringz.FirstNPartsX("", 2, "<|>")
				Ω(x).To(Equal(""))
				Ω(remain).To(Equal(""))
			})
			Convey("LastNParts", func() {
				x, remain := stringz.LastNPartsX("", 2, "<|>")
				Ω(x).To(Equal(""))
				Ω(remain).To(Equal(""))
			})
		})

		input := "one<|>two<|>three<|>four"
		Convey("FirstNParts", func() {
			Convey("0", func() {
				x, remain := stringz.FirstNPartsX(input, 0, "<|>")
				Ω(x).To(Equal(""))
				Ω(remain).To(Equal("one<|>two<|>three<|>four"))
			})
			Convey("2", func() {
				x, remain := stringz.FirstNPartsX(input, 2, "<|>")
				Ω(x).To(Equal("one<|>two"))
				Ω(remain).To(Equal("three<|>four"))
			})
			Convey("4", func() {
				x, remain := stringz.FirstNPartsX(input, 4, "<|>")
				Ω(x).To(Equal("one<|>two<|>three<|>four"))
				Ω(remain).To(Equal(""))
			})
			Convey("10", func() {
				x, remain := stringz.FirstNPartsX(input, 10, "<|>")
				Ω(x).To(Equal("one<|>two<|>three<|>four"))
				Ω(remain).To(Equal(""))
			})
		})
		Convey("LastNParts", func() {
			Convey("0", func() {
				x, remain := stringz.LastNPartsX(input, 0, "<|>")
				Ω(x).To(Equal(""))
				Ω(remain).To(Equal("one<|>two<|>three<|>four"))
			})
			Convey("1", func() {
				x, remain := stringz.LastNPartsX(input, 1, "<|>")
				Ω(x).To(Equal("four"))
				Ω(remain).To(Equal("one<|>two<|>three"))
			})
			Convey("3", func() {
				x, remain := stringz.LastNPartsX(input, 3, "<|>")
				Ω(x).To(Equal("two<|>three<|>four"))
				Ω(remain).To(Equal("one"))
			})
			Convey("4", func() {
				x, remain := stringz.LastNPartsX(input, 4, "<|>")
				Ω(x).To(Equal("one<|>two<|>three<|>four"))
				Ω(remain).To(Equal(""))
			})
			Convey("10", func() {
				x, remain := stringz.LastNPartsX(input, 10, "<|>")
				Ω(x).To(Equal("one<|>two<|>three<|>four"))
				Ω(remain).To(Equal(""))
			})
		})
	})
}
