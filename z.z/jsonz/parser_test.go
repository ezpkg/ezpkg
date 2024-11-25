package jsonz_test

import (
	"testing"

	. "ezpkg.io/conveyz"
	"ezpkg.io/jsonz"
	"ezpkg.io/stringz"
	. "ezpkg.io/testingz"
)

func TestParse(t *testing.T) {
	ΩxNoDiff := ΩxNoDiffByLineZ

	Convey("Parse", t, func() {
		parse := func(in string) (_ string) {
			var b stringz.Builder
			for item, err := range jsonz.Parse([]byte(in)) {
				if err != nil {
					b.Printf("[ERROR] %v", err)
				} else {
					b.Printf("%+v\n", item)
				}
			}
			return b.String()
		}
		Convey("simple", func() {
			Convey("number", func() {
				s := parse("1234")
				ΩxNoDiff(s, `→ 1234`)
			})
			Convey("string", func() {
				s := parse(`"foo"`)
				ΩxNoDiff(s, `→ "foo"`)
			})
			Convey("null", func() {
				s := parse(`null`)
				ΩxNoDiff(s, `→ null`)
			})
			Convey("true", func() {
				s := parse(`true`)
				ΩxNoDiff(s, `→ true`)
			})
			Convey("false", func() {
				s := parse(`false`)
				ΩxNoDiff(s, `→ false`)
			})
			Convey("array", func() {
				s := parse(`[1,"2",3]`)
				ΩxNoDiff(s, `
[0] → 1
[1] → "2"
[2] → 3`)
			})
			Convey("object", func() {
				s := parse(`{"a":1,"b":"2","c":3}`)
				ΩxNoDiff(s, `
."a" → 1
."b" → "2"
."c" → 3`)
			})
		})
	})
}
