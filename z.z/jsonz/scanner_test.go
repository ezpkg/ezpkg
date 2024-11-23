package draft

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	jtest "ezpkg.io/-/jsonz_test"
)

func TestNextToken(t *testing.T) {
	for _, test := range jtest.LargeSet {
		t.Run(test.Name, func(t *testing.T) {
			buf := make([]byte, 0, len(test.Data))
			w := bytes.NewBuffer(buf)
			last := test.Data
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
