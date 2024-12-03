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
		tcase := jtest.GetTestcase("pass01.json")
		out, err := jsonz.Reconstruct(tcase.Data)
		Ω(err).ToNot(HaveOccurred())

		fmt.Printf("--- reconstruct ---\n%s\n---\n", out)

		actual := reformatWithStdjson(out)
		expect := reformatWithStdjson(tcase.Data)
		ΩxNoDiffByLine(actual, expect)
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
