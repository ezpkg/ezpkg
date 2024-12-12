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
		// 🐝Example: minify json
		b := jsoniter.NewBuilder("", "")
		for item, err := range jsoniter.Parse(data) {
			errorz.MustZ(err)
			b.Add(item.Key, item.Token)
		}
		out := errorz.Must(b.Bytes())
		fmt.Printf("\n--- minify ---\n%s\n----------\n", out)
	}
	{
		// 🦋Example: reformat json
		b := jsoniter.NewBuilder("👉   ", "\t")
		for item, err := range jsoniter.Parse(data) {
			errorz.MustZ(err)
			b.Add(item.Key, item.Token)
		}
		out := errorz.Must(b.Bytes())
		fmt.Printf("\n--- reformat ---\n%s\n----------\n", out)
	}
}
