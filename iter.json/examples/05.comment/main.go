package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"ezpkg.io/errorz"
	iterjson "ezpkg.io/iter.json"
)

func main() {
	_, file, _, _ := runtime.Caller(0)
	data := errorz.Must(os.ReadFile(filepath.Dir(file) + "/../alice.json"))
	{
		// 🐶Example: add comment with line number
		i, newlineIdx, maxIdx := 0, 0, 50
		b := iterjson.NewBuilder("", "    ")
		for item, err := range iterjson.Parse(data) {
			errorz.MustZ(err)
			b.WriteComma(item.Token.Type())

			// 👉 add comment
			if i > 0 {
				length := b.Len() - newlineIdx
				fmt.Fprint(b, strings.Repeat(" ", maxIdx-length))
				fmt.Fprintf(b, "// %2d", i)
			}
			i++

			b.WriteNewline(item.Token.Type())
			newlineIdx = b.Len() // save the newline index

			b.Add(item.Key, item.Token)
		}
		length := b.Len() - newlineIdx
		fmt.Fprint(b, strings.Repeat(" ", maxIdx-length))
		fmt.Fprintf(b, "// %2d", i)

		out := errorz.Must(b.Bytes())
		fmt.Printf("\n--- comment ---\n%s\n----------\n", out)
	}
}
