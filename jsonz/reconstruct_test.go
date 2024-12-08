package jsonz_test

import (
	stdjson "encoding/json"
	"fmt"
	"testing"

	jtest "ezpkg.io/-/jsonz_test"
	. "ezpkg.io/conveyz"
	"ezpkg.io/jsonz"
	. "ezpkg.io/testingz"
)

func TestReconstruct(t *testing.T) {
	Convey("Reconstruct", t, func() {
		Convey("no ident", func() {
			tcase := jtest.GetTestcase("pass01.json")
			out, err := jsonz.Reconstruct(tcase.Data)
			Ω(err).ToNot(HaveOccurred())

			fmt.Printf("\n--- reconstruct ---\n%s\n---\n", out)

			actual := reformatWithStdjson(out)
			expect := reformatWithStdjson(tcase.Data)
			ΩxNoDiffByLine(actual, expect)
		})
		Convey("with ident", func() {
			tcase := jtest.GetTestcase("pass01.json")
			out, err := jsonz.Reformat(tcase.Data, "", "\t")
			Ω(err).ToNot(HaveOccurred())

			fmt.Printf("\n--- reformat ---\n%s\n---\n", out)
			ΩxNoDiffByLine(string(out), tcase.ExpectFormat)
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
