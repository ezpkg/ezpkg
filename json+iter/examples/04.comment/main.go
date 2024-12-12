package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"ezpkg.io/errorz"
	jsoniter "ezpkg.io/json+iter"
)

func main() {
	_, file, _, _ := runtime.Caller(0)
	data := errorz.Must(os.ReadFile(filepath.Dir(file) + "/../alice.json"))
	{
		// ??Example: add comment with line number
		i, newlineIdx, maxIdx := 0, 0, 50
		b := jsoniter.NewBuilder("", "    ")
		for item, err := range jsoniter.Parse(data) {
			errorz.MustZ(err)
			b.WriteComma(item.Token.Type())

			// 👉 add comment
			if i > 0 {
				length := b.Len() - newlineIdx
				writeSpace(b, maxIdx-length)
				fprintf(b, "// %2d", i)
			}
			i++

			b.WriteNewline(item.Token.Type())
			newlineIdx = b.Len() // save the newline index

			b.AddRaw(item.Key, item.Token)
		}
		length := b.Len() - newlineIdx
		writeSpace(b, maxIdx-length)
		fprintf(b, "// %2d", i)

		out := errorz.Must(b.Bytes())
		fmt.Printf("\n--- comment ---\n%s\n----------\n", out)
	}
}

var space = []byte(" ")

func writeSpace(w io.Writer, n int) {
	for i := 0; i < n; i++ {
		errorz.Must(w.Write(space))
	}
}

func fprintf(w io.Writer, format string, args ...any) {
	errorz.Must(fmt.Fprintf(w, format, args...))
}
