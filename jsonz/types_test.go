package jsonz_test

import (
	"testing"

	. "ezpkg.io/conveyz"
	. "ezpkg.io/jsonz"
)

func TestTypes(t *testing.T) {
	Convey("Types", t, func() {
		Convey("PathItem", func() {
			Convey("array ✅", func() {
				path, err := NewPathItem(1)
				Ω(err).ToNot(HaveOccurred())
				Ω(path).To(Equal(PathItem{
					Index: 1,
					Token: TokenArrayOpen.New(),
				}))
			})
			Convey("object: unquote ✅", func() {
				path, err := NewPathItem("foo")
				Ω(err).ToNot(HaveOccurred())
				Ω(path).To(Equal(PathItem{
					Key:   MustRawToken([]byte(`"foo"`)),
					Token: TokenObjectOpen.New(),
				}))
			})
			Convey("object: quote ✅", func() {
				path, err := NewPathItem(`"foo"`)
				Ω(err).ToNot(HaveOccurred())
				Ω(path).To(Equal(PathItem{
					Key:   MustRawToken([]byte(`"foo"`)),
					Token: TokenObjectOpen.New(),
				}))
			})
			Convey("object: unquote → quote ✅", func() {
				path, err := NewPathItem(`a"b`)
				Ω(err).ToNot(HaveOccurred())
				Ω(path).To(Equal(PathItem{
					Key:   MustRawToken([]byte(`"a\"b"`)),
					Token: TokenObjectOpen.New(),
				}))
			})
			Convey("array: neg ❌", func() {
				_, err := NewPathItem(-1)
				Ω(err).To(HaveOccurred())
			})
			Convey("array: float ❌", func() {
				_, err := NewPathItem(0.2)
				Ω(err).To(HaveOccurred())
			})
			Convey("object: unquote ❌", func() {
				_, err := NewPathItem(`"a\"b`)
				Ω(err).To(HaveOccurred())
			})
		})
		Convey("RawPath", func() {
			path, err := NewRawPath("foo", 1, "bar", 42)
			Ω(err).ToNot(HaveOccurred())
			Ω(path).To(Equal(RawPath{
				PathItem{
					Key:   MustRawToken([]byte(`"foo"`)),
					Token: TokenObjectOpen.New(),
				},
				PathItem{
					Index: 1,
					Token: TokenArrayOpen.New(),
				},
				PathItem{
					Key:   MustRawToken([]byte(`"bar"`)),
					Token: TokenObjectOpen.New(),
				},
				PathItem{
					Index: 42,
					Token: TokenArrayOpen.New(),
				},
			}))
		})
		Convey("RawPath: compare", func() {
			path1, _ := NewRawPath("foo", 1, "bar", 42)
			path2, _ := NewRawPath("foo", 1, "bar", 42)
			path3, _ := NewRawPath("foo", 1, "bar", 43)

			Convey("Equal", func() {
				Convey("RawPath", func() {
					Ω(path1.Match(path2)).To(BeTrue())
					Ω(path1.Match(path3)).To(BeFalse())
				})
				Convey("...any", func() {
					Ω(path1.Match("foo", 1, "bar", 42)).To(BeTrue())
					Ω(path1.Match("foo", 1, "bar", 43)).To(BeFalse())
					Ω(path1.Match("foo", 1, "bar")).To(BeFalse())
				})
				Convey("[]any", func() {
					slice := []any{"foo", 1, "bar", 42}
					Ω(path1.Match(slice)).To(BeTrue())
				})
			})
		})
	})
}
