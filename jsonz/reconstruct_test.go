package jsonz_test

import (
	"bytes"
	stdjson "encoding/json"
	"fmt"
	"math"
	"testing"

	. "ezpkg.io/conveyz"
	. "ezpkg.io/jsonz"
	"ezpkg.io/jsonz/test"
	. "ezpkg.io/testingz"
)

func TestReconstruct(t *testing.T) {
	Convey("Reconstruct", t, func() {
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
	Convey("Builder", t, func() {
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
