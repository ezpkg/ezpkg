package draft

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	jtest "ezpkg.io/-/jsonz_test"
)

func TestNextToken(t *testing.T) {
	run := func(tcase jtest.Testcase) {
		t.Run(tcase.Name, func(t *testing.T) {
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
				must(fmt.Fprintf(w, "%s", token))
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
}

func TestScan(t *testing.T) {
	run := func(tcase jtest.Testcase) {
		t.Run(tcase.Name, func(t *testing.T) {
			buf := make([]byte, 0, len(tcase.Data))
			w := bytes.NewBuffer(buf)
			for token, err := range Scan(tcase.Data) {
				if err != nil {
					t.Errorf("error: %v", err)
					return
				}
				w.Write(token)
			}

			var v any
			must(0, json.Unmarshal(w.Bytes(), &v))
		})
	}

	run(jtest.GetTestcase("pass01.json"))
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func panicf(format string, args ...any) {
	panic(fmt.Sprintf(format, args...))
}
