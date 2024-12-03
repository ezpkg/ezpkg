package jsonz_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	jtest "ezpkg.io/-/jsonz_test"
	. "ezpkg.io/conveyz"
	. "ezpkg.io/jsonz"
	"ezpkg.io/testingz"
)

func TestNextToken(t *testing.T) {
	Convey("NextToken", t, func() {
		run := func(tcase jtest.Testcase) {
			Convey(tcase.Name, func() {
				buf := make([]byte, 0, len(tcase.Data))
				w := bytes.NewBuffer(buf)
				last := tcase.Data
				for {
					token, remain, err := NextToken(last)
					if err != nil {
						t.Errorf("error: %v", err)
						return
					}
					if len(remain) >= len(last) {
						panicf("unexpected: remain[%v] >= last[%v]", len(remain), len(last))
					}
					must(w.Write(token.Bytes()))
					if len(remain) == 0 {
						break
					}
					last = remain
				}

				var v any
				must(0, json.Unmarshal(w.Bytes(), &v))
			})
		}

		run(jtest.GetTestcase("pass01.json"))
		for _, test := range jtest.LargeSet {
			run(test)
		}
	})
}

func TestScan(t *testing.T) {
	Convey("Scan", t, func() {
		run := func(tcase jtest.Testcase) {
			Convey(tcase.Name, func() {
				_number0, _number1, _string0 := false, false, false

				buf := make([]byte, 0, len(tcase.Data))
				w := bytes.NewBuffer(buf)
				for token, err := range Scan(tcase.Data) {
					if err != nil {
						t.Errorf("error: %v", err)
						return
					}
					must(fmt.Fprintln(w, token))

					// verify a few tokens
					x, err := token.ParseValue()
					switch token.String() {
					case "null":
						assert(t, err == nil && x == nil, "expected nil")
						assert(t, token.Type() == TokenNull, "expected null")

					case "true":
						assert(t, err == nil && x == true, "expected true")
						assert(t, token.Type() == TokenTrue, "expected true")

					case "false":
						assert(t, err == nil && x == false, "expected false")
						assert(t, token.Type() == TokenFalse, "expected false")

					case "1234567890":
						_number0 = true
						assert(t, token.Type() == TokenNumber, "expected number")
						n, err := token.ParseNumber()
						assert(t, err == nil && n == float64(1234567890), "expected 1234567890")

					case "1.234567890E+34":
						_number1 = true
						assert(t, token.Type() == TokenNumber, "expected number")
						n, err := token.ParseNumber()
						assert(t, err == nil && n == float64(1.234567890e+34), "expected 1.234567890E+34")

					case `"\u0123\u4567\u89AB\uCDEF\uabcd\uef4A"`:
						_string0 = true
						assert(t, token.Type() == TokenString, "expected string")
						s, err := token.ParseString()
						assert(t, err == nil && s == "\u0123\u4567\u89AB\uCDEF\uabcd\uef4A", "expected \u0123\u4567\u89AB\uCDEF\uabcd\uef4A")
					}
				}

				Ω(tcase.ExpectTokens).ToNot(BeEmpty())
				testingz.ΩxNoDiffByLine(w.String(), tcase.ExpectTokens)

				var v any
				must(0, json.Unmarshal(w.Bytes(), &v))

				if !(_number0 && _number1 && _string0) {
					t.Errorf("expected number0(%v), number1(%v), string0(%v)", _number0, _number1, _string0)
				}
			})
		}

		run(jtest.GetTestcase("pass01.json"))
	})
}

func assertToken(t *testing.T, token RawToken, typ TokenType, value any) {
	switch token.Type() {
	case TokenNumber:
		n, err := token.ParseNumber()
		assert(t, err == nil && n == value.(float64), "expected %v", value)
	case TokenString:
		s, err := token.ParseString()
		assert(t, err == nil && s == value.(string), "expected %v", value)
	default:
		panic("unreachable")
	}
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func assert(t *testing.T, cond bool, msg string, args ...any) {
	if !cond {
		t.Errorf("❌ "+msg, args...)
	}
}

func panicf(format string, args ...any) {
	panic(fmt.Sprintf(format, args...))
}
