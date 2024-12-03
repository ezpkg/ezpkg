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
		parse := func(in string) (string, error) {
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
					return b.String(), err
				}
				_, err0 := item.ParseValue()
				if err0 != nil {
					pr("[ERROR] %v\n", err0)
					return b.String(), err0
				}
				pr("%+v\n", item)
			}
			return b.String(), nil
		}
		Convey("simple", func() {
			Convey("number", func() {
				s, _ := parse("1234")
				ΩxNoDiff(s, `→ 1234`)
			})
			Convey("string", func() {
				s, _ := parse(`"foo"`)
				ΩxNoDiff(s, `→ "foo"`)
			})
			Convey("null", func() {
				s, _ := parse(`null`)
				ΩxNoDiff(s, `→ null`)
			})
			Convey("true", func() {
				s, _ := parse(`true`)
				ΩxNoDiff(s, `→ true`)
			})
			Convey("false", func() {
				s, _ := parse(`false`)
				ΩxNoDiff(s, `→ false`)
			})
			Convey("empty array", func() {
				s, _ := parse(`[]`)
				ΩxNoDiff(s, `
→ [
→ ]`)
			})
			Convey("array", func() {
				s, _ := parse(`[1,"2",3]`)
				ΩxNoDiff(s, `
    → [
[0] → 1
[1] → "2"
[2] → 3
    → ]`)
			})
			Convey("empty object", func() {
				s, _ := parse(`{}`)
				ΩxNoDiff(s, `
→ {
→ }`)
			})
			Convey("object", func() {
				s, _ := parse(`{"a":1,"b":"2","c":3}`)
				ΩxNoDiff(s, `
     → {
."a" → 1
."b" → "2"
."c" → 3
 → }`)
			})
		})
		Convey("nested", func() {
			Convey("array 2x2", func() {
				s, _ := parse(`[[1,2],[3,4]]`)
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
			Convey("array empty", func() {
				s, _ := parse(`[[]]`)
				ΩxNoDiff(s, `
    → [
[0] → [
[0] → ]
    → ]`)
			})
		})
		Convey("pass01.json", func() {
			tcase := jtest.GetTestcase("pass01.json")
			s, _ := parse(string(tcase.Data))
			ΩxNoDiff(s, tcase.ExpectParse)
		})
		Convey("failures", func() {
			for _, tcase := range jtest.SimpleSet {
				if !tcase.Bad {
					continue
				}
				Convey(tcase.Name, func() {
					_, err := parse(string(tcase.Data))
					fmt.Printf("%s:\t%v\n", tcase.Name, err)
					Ω(err).Should(HaveOccurred())
				})
			}
		})
	})
}
