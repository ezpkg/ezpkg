package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"ezpkg.io/errorz"
	"ezpkg.io/jsonz"
)

func main() {
	_, file, _, _ := runtime.Caller(0)
	data := errorz.Must(os.ReadFile(filepath.Dir(file) + "/../alice.json"))
	{
		// 🐞Example: print with line number
		i := 0
		b := jsonz.NewBuilder("", "    ")
		for item, err := range jsonz.Parse(data) {
			i++
			errorz.MustZ(err)
			b.WriteNewline(item.Token.Type())

			// 👉 add line number
			fprintf(b, "%3d    ", i)
			b.AddRaw(item.Key, item.Token)
		}
		out := errorz.Must(b.Bytes())
		fmt.Printf("\n--- line number ---\n%s\n----------\n", out)
	}
}

func fprintf(w io.Writer, format string, args ...any) {
	errorz.Must(fmt.Fprintf(w, format, args...))
}
