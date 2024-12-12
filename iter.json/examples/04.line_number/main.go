package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"ezpkg.io/errorz"
	iterjson "ezpkg.io/iter.json"
)

func main() {
	_, file, _, _ := runtime.Caller(0)
	data := errorz.Must(os.ReadFile(filepath.Dir(file) + "/../alice.nguyen.json"))
	{
		// 🐞Example: print with line number
		i := 0
		b := iterjson.NewBuilder("", "    ")
		for item, err := range iterjson.Parse(data) {
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
