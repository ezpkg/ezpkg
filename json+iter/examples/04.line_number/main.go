package main

import (
	"fmt"
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
		// 🐞Example: print with line number
		i := 0
		b := jsoniter.NewBuilder("", "    ")
		for item, err := range jsoniter.Parse(data) {
			i++
			errorz.MustZ(err)
			b.WriteNewline(item.Token.Type())

			// 👉 add line number
			fmt.Fprintf(b, "%3d    ", i)
			b.Add(item.Key, item.Token)
		}
		out := errorz.Must(b.Bytes())
		fmt.Printf("\n--- line number ---\n%s\n----------\n", out)
	}
}
