package jsonz_test

import (
	"bytes"
	stdjson "encoding/json"
	"fmt"
	"testing"

	. "ezpkg.io/conveyz"
	"ezpkg.io/jsonz"
	"ezpkg.io/jsonz/test"
	. "ezpkg.io/testingz"
)

func TestReconstruct(t *testing.T) {
	Convey("Reconstruct", t, func() {
		Convey("no indent", func() {
			tcase := test.GetTestcase("pass01.json")
			out, err := jsonz.Reconstruct(tcase.Data)
			Ω(err).ToNot(HaveOccurred())

			fmt.Printf("\n--- reconstruct ---\n%s\n---\n", out)

			actual := reformatWithStdjson(out)
			expect := reformatWithStdjson(tcase.Data)
			ΩxNoDiffByLine(actual, expect)
		})
		Convey("with indent", func() {
			tcase := test.GetTestcase("pass01.json")
			out, err := jsonz.Reformat(tcase.Data, "→ ", "\t")
			Ω(err).ToNot(HaveOccurred())

			fmt.Printf("\n--- reformat ---\n%s\n---\n", out)
			ΩxNoDiffByLine(string(out), tcase.ExpectFormat)
		})
	})
	Convey("Builder", t, func() {
		Convey("no indent", func() {
			tcase := test.GetTestcase("pass01.json")
			b := jsonz.NewBuilder("", "")
			for item, err := range jsonz.Parse(tcase.Data) {
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
			b := jsonz.NewBuilder("→ ", "\t")
			for item, err := range jsonz.Parse(tcase.Data) {
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
