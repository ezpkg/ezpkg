package draft

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/iOliverNguyen/jsonz/testingz"
)

func TestNextToken(t *testing.T) {
	for _, test := range ValidTestcases() {
		t.Run(test.File, func(t *testing.T) {
			buf := make([]byte, 0, len(test.Data))
			w := bytes.NewBuffer(buf)
			for {
				token, remain, err := NextToken(test.DataStr())
				if err != nil {
					t.Errorf("error: %v", err)
					return
				}
				if remain == "" {
					break
				}
				Must(fmt.Fprintf(w, token))
			}

			var v any
			Must(0, json.Unmarshal(w.Bytes(), &v))
		})
	}
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
