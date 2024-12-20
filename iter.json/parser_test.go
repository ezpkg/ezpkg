package iterjson_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	. "ezpkg.io/conveyz"
	iterjson "ezpkg.io/iter.json"
	"ezpkg.io/iter.json/test"
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
			for item, err := range iterjson.Parse([]byte(in)) {
				if err != nil {
					pr("[ERROR] %v\n", err)
					return b.String(), err
				}
				if item.Token.IsValue() {
					_, err0 := item.GetValue()
					if err0 != nil {
						pr("[ERROR] %v\n", err0)
						return b.String(), err0
					}
				}
				pr("L%v  %+v\t\n", item.Level, item)
			}
			return b.String(), nil
		}
		Convey("simple", func() {
			Convey("number", func() {
				s, _ := parse("1234")
				ΩxNoDiff(s, `L0 → 1234`)
			})
			Convey("string", func() {
				s, _ := parse(`"foo"`)
				ΩxNoDiff(s, `L0 → "foo"`)
			})
			Convey("null", func() {
				s, _ := parse(`null`)
				ΩxNoDiff(s, `L0 → null`)
			})
			Convey("true", func() {
				s, _ := parse(`true`)
				ΩxNoDiff(s, `L0 → true`)
			})
			Convey("false", func() {
				s, _ := parse(`false`)
				ΩxNoDiff(s, `L0 → false`)
			})
			Convey("empty array", func() {
				s, _ := parse(`[]`)
				ΩxNoDiff(s, `
L0 → [
L0 → ]`)
			})
			Convey("array", func() {
				s, _ := parse(`[1,"2",3]`)
				ΩxNoDiff(s, `
L0     → [
L1  0  → 1
L1  1  → "2"
L1  2  → 3
L0     → ]`)
			})
			Convey("empty object", func() {
				s, _ := parse(`{}`)
				ΩxNoDiff(s, `
L0  → {
L0  → }
`)
			})
			Convey("object", func() {
				s, _ := parse(`{"a":1,"b":"2","c":3}`)
				ΩxNoDiff(s, `
L0    → {
L1  a → 1
L1  b → "2"
L1  c → 3
L0    → }
`)
			})
		})
		Convey("nested", func() {
			Convey("array 2x2", func() {
				s, _ := parse(`[[1,2],[3,4]]`)
				ΩxNoDiff(s, `
L0       → [
L1  0    → [
L2  0.0  → 1
L2  0.1  → 2
L1  0    → ]
L1  1    → [
L2  1.0  → 3
L2  1.1  → 4
L1  1    → ]
L0       → ]
`)
			})
			Convey("array empty", func() {
				s, _ := parse(`[[]]`)
				ΩxNoDiff(s, `
L0     → [
L1  0  → [
L1  0  → ]
L0     → ]
`)
			})
		})
		Convey("GetPath", func() {
			input := `{"name": "Alice", "age": 24, "scores": [9, {"math": 10}, 8], "address": {"city": "The Sun", "zip": 10101}}`

			var b stringz.Builder
			for item, err := range iterjson.Parse([]byte(input)) {
				if err != nil {
					b.Printf("[ERROR] %v\n", err)
				} else {
					b.Printf("%v → %v\n", item.GetPath(), item.Token)
					b.Printf("%v → %v\n", item.GetPathString(), item.Token)
					b.Printf("%v → %v\n", item.GetAltPathString(), item.Token)
					b.Printf("%v → %v\n", item.GetRawPath(), item.Token)
					b.Println()
				}
			}
			ΩxNoDiff(b.String(), `
 → {
 → {
 → {
 → {

name → "Alice"
name → "Alice"
.name → "Alice"
name → "Alice"

age → 24
age → 24
.age → 24
age → 24

scores → [
scores → [
.scores → [
scores → [

scores.0 → 9
scores.0 → 9
.scores[0] → 9
scores.0 → 9

scores.1 → {
scores.1 → {
.scores[1] → {
scores.1 → {

scores.1.math → 10
scores.1.math → 10
.scores[1].math → 10
scores.1.math → 10

scores.1 → }
scores.1 → }
.scores[1] → }
scores.1 → }

scores.2 → 8
scores.2 → 8
.scores[2] → 8
scores.2 → 8

scores → ]
scores → ]
.scores → ]
scores → ]

address → {
address → {
.address → {
address → {

address.city → "The Sun"
address.city → "The Sun"
.address.city → "The Sun"
address.city → "The Sun"

address.zip → 10101
address.zip → 10101
.address.zip → 10101
address.zip → 10101

address → }
address → }
.address → }
address → }

 → }
 → }
 → }
 → }
`)
		})
		Convey("pass01.json", func() {
			tcase := test.GetTestcase("pass01.json")
			s, _ := parse(string(tcase.Data))
			ΩxNoDiff(s, tcase.ExpectParse)
		})
		Convey("failures", func() {
			for _, tcase := range test.SimpleSet {
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
