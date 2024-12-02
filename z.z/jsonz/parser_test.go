package jsonz_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	jtest "ezpkg.io/-/jsonz_test"
	. "ezpkg.io/conveyz"
	"ezpkg.io/jsonz"
	"ezpkg.io/stringz"
	. "ezpkg.io/testingz"
)

func TestParse(t *testing.T) {
	ΩxNoDiff := ΩxNoDiffByLineZ

	debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))
	Convey("Parse", t, func() {
		parse := func(in string) (_ string) {
			var b stringz.Builder
			pr := func(msg string, args ...any) {
				if debug {
					fmt.Printf(msg, args...)
				}
				b.Printf(msg, args...)
			}
			for item, err := range jsonz.Parse([]byte(in)) {
				if err != nil {
					pr("[ERROR] %v\n", err)
				} else {
					pr("%+v\n", item)
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
			Convey("empty array", func() {
				s := parse(`[]`)
				ΩxNoDiff(s, `
→ [
→ ]`)
			})
			Convey("array", func() {
				s := parse(`[1,"2",3]`)
				ΩxNoDiff(s, `
    → [
[0] → 1
[1] → "2"
[2] → 3
    → ]`)
			})
			Convey("empty object", func() {
				s := parse(`{}`)
				ΩxNoDiff(s, `
→ {
→ }`)
			})
			Convey("object", func() {
				s := parse(`{"a":1,"b":"2","c":3}`)
				ΩxNoDiff(s, `
     → {
."a" → 1
."b" → "2"
."c" → 3
 → }`)
			})
		})
		Convey("nested array", func() {
			Convey("2x2", func() {
				s := parse(`[[1,2],[3,4]]`)
				ΩxNoDiff(s, `
       → [
[0]    → [
[0][0] → 1
[0][1] → 2
[0]    → ]
[1]    → [
[1][0] → 3
[1][1] → 4
[1]    → ]
       → ]`)
			})
			Convey("empty", func() {
				s := parse(`[[]]`)
				ΩxNoDiff(s, `
[0] → [0] → []`)
			})
		})
		Convey("pass01.json", func() {
			tcase := jtest.GetTestcase("pass01.json")
			s := parse(string(tcase.Data))
			ΩxNoDiff(s, tcase.ExpectParse)
		})
	})
}
