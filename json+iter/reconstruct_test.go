package jsoniter_test

import (
	"bytes"
	stdjson "encoding/json"
	"fmt"
	"math"
	"testing"

	. "ezpkg.io/conveyz"
	. "ezpkg.io/json+iter"
	"ezpkg.io/json+iter/test"
	. "ezpkg.io/testingz"
)

func TestReconstruct(t *testing.T) {
	Convey("Reconstruct", t, func() {
		Convey("Reformat", func() {
			Convey("no indent", func() {
				tcase := test.GetTestcase("pass01.json")
				out, err := Reconstruct(tcase.Data)
				Ω(err).ToNot(HaveOccurred())

				fmt.Printf("\n--- reconstruct ---\n%s\n---\n", out)

				actual := reformatWithStdjson(out)
				expect := reformatWithStdjson(tcase.Data)
				ΩxNoDiffByLine(actual, expect)
			})
			Convey("with indent", func() {
				tcase := test.GetTestcase("pass01.json")
				out, err := Reformat(tcase.Data, "→ ", "\t")
				Ω(err).ToNot(HaveOccurred())

				fmt.Printf("\n--- reformat ---\n%s\n---\n", out)
				ΩxNoDiffByLine(string(out), tcase.ExpectFormat)
			})
		})
		Convey("Builder", func() {
			Convey("no indent", func() {
				tcase := test.GetTestcase("pass01.json")
				b := NewBuilder("", "")
				for item, err := range Parse(tcase.Data) {
					Ω(err).ToNot(HaveOccurred())
					b.AddRaw(item.Key, item.Token)
				}
				out, err := b.Bytes()
				Ω(err).ToNot(HaveOccurred())

				fmt.Printf("\n--- reconstruct ---\n%s\n---\n", out)

				// verify that the output json is correct
				actual := reformatWithStdjson(out)
				expect := reformatWithStdjson(tcase.Data)
				ΩxNoDiffByLine(actual, expect)
			})
			Convey("with indent", func() {
				tcase := test.GetTestcase("pass01.json")
				b := NewBuilder("→ ", "\t")
				for item, err := range Parse(tcase.Data) {
					Ω(err).ToNot(HaveOccurred())
					b.AddRaw(item.Key, item.Token)
				}
				out, err := b.Bytes()
				Ω(err).ToNot(HaveOccurred())

				fmt.Printf("\n--- reformat ---\n%s\n---\n", out)

				ΩxNoDiffByLine(string(out), tcase.ExpectFormat)

				// verify that the output json is correct
				out = bytes.ReplaceAll(out, []byte("→ "), []byte(""))
				actual := reformatWithStdjson(out)
				expect := reformatWithStdjson(tcase.Data)
				ΩxNoDiffByLine(actual, expect)
			})
			Convey("Add+AddRaw", func() {
				expected := `{
  "key1": 123,
  "key2": [
    0.42,
    1,
    "2",
    3,
    "four",
    true,
    0,
    0
  ]
}`
				Convey("AddRaw()", func() {
					b := NewBuilder("", "  ")
					b.AddRaw(RawToken{}, TokenObjectOpen.New())
					b.AddRaw(StringToken("key1"), NumberToken(123))
					b.AddRaw(StringToken("key2"), TokenArrayOpen.New())
					b.AddRaw(RawToken{}, NumberToken(0.42))
					b.AddRaw(RawToken{}, IntToken(1))
					b.AddRaw(RawToken{}, StringToken("2"))
					b.AddRaw(RawToken{}, MustRawToken([]byte("3")))
					b.AddRaw(RawToken{}, MustRawToken([]byte(`"four"`)))
					b.AddRaw(RawToken{}, BoolToken(true))
					b.AddRaw(RawToken{}, NumberToken(math.NaN()))   // fallback to 0
					b.AddRaw(RawToken{}, NumberToken(math.Inf(-1))) // fallback to 0
					b.AddRaw(RawToken{}, TokenArrayClose.New())
					b.AddRaw(RawToken{}, TokenObjectClose.New())

					out := string(must(b.Bytes()))
					ΩxNoDiffByLine(out, expected)
				})
				Convey("Add() - RawToken", func() {
					b := NewBuilder("", "  ")
					b.Add("", TokenObjectOpen)
					b.Add("key1", 123)
					b.Add("key2", TokenArrayOpen)
					b.Add("", NumberToken(0.42))
					b.Add("", IntToken(1))
					b.Add("", StringToken("2"))
					b.Add("", MustRawToken([]byte("3")))
					b.Add("", MustRawToken([]byte(`"four"`)))
					b.Add("", BoolToken(true))
					b.Add("", NumberToken(math.NaN()))   // fallback to 0
					b.Add("", NumberToken(math.Inf(-1))) // fallback to 0
					b.Add("", TokenArrayClose.New())
					b.Add("", TokenObjectClose.New())

					out := string(must(b.Bytes()))
					ΩxNoDiffByLine(out, expected)
				})
				Convey("Add() - any", func() {
					b := NewBuilder("", "  ")
					b.Add("", TokenObjectOpen)
					b.Add("key1", 123)
					b.Add("key2", TokenArrayOpen)
					b.Add("", 0.42)
					b.Add("", 1)
					b.Add("", "2")
					b.Add("", []byte("3"))
					b.Add("", []byte(`"four"`))
					b.Add("", true)
					b.Add("", math.NaN())   // fallback to 0
					b.Add("", math.Inf(-1)) // fallback to 0
					b.Add("", TokenArrayClose)
					b.Add("", TokenObjectClose)

					out := string(must(b.Bytes()))
					ΩxNoDiffByLine(out, expected)
				})
				Convey("Add() - fallback", func() {
					b := NewBuilder("", "  ")
					b.Add("", TokenObjectOpen)
					b.Add("key1", 123)
					b.Add("key2", []any{0.42, 1, "2", 3, "four", true})
					b.Add("key3", "ok")
					b.Add("", TokenObjectClose)

					out := string(must(b.Bytes()))
					ΩxNoDiffByChar(out, `{
  "key1": 123,
  "key2": [0.42,1,"2",3,"four",true],
  "key3": "ok"
}`)
				})
			})
			Convey("SkipEmptyStructures", func() {
				exec := func(in string) string {
					b := NewBuilder("", "")
					b.SetSkipEmptyStructures(true)
					for item, err := range Parse([]byte(in)) {
						Ω(err).ToNot(HaveOccurred())
						b.AddRaw(item.Key, item.Token)
					}
					return string(must(b.Bytes()))
				}

				Convey("[]", func() {
					out := exec(`[]`)
					Ω(out).To(Equal("null"))
				})
				Convey("{}", func() {
					out := exec(`{}`)
					Ω(out).To(Equal("null"))
				})
				Convey(`[1,[],{},2]`, func() {
					out := exec(`[1,[],{},2]`)
					Ω(out).To(Equal("[1,2]"))
				})
				Convey(`[[[[],1,[]]]]`, func() {
					out := exec(`[[[[],1,[]]]]`)
					Ω(out).To(Equal("[[[1]]]"))
				})
				Convey(`[[],"foo",{}]`, func() {
					out := exec(`[[],"foo",{}]`)
					Ω(out).To(Equal(`["foo"]`))
				})
				Convey(`{"a":{},"b":42,"c":[]}`, func() {
					out := exec(`{"a":{},"b":42,"c":[]}`)
					Ω(out).To(Equal(`{"b":42}`))
				})
				Convey(`complex`, func() {
					out := exec(`
{
  "a": {},
  "_": [
    [],
    {
      "d": {},
      "e": "1",
      "f": [[1]]
    }
  ],
  "c": []
}`)
					Ω(out).To(Equal(`{"_":[{"e":"1","f":[[1]]}]}`))
				})
			})
		})
	})
}

func reformatWithStdjson(in []byte) string {
	var x any
	err := stdjson.Unmarshal(in, &x)
	Ω(err).ToNot(HaveOccurred())

	out, err := stdjson.MarshalIndent(x, "", "  ")
	Ω(err).ToNot(HaveOccurred())
	return string(out)
}
