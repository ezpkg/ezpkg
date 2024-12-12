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
	data := errorz.Must(os.ReadFile(filepath.Dir(file) + "/../alice.json"))

	{
		// 🐝Example: minify json
		b := iterjson.NewBuilder("", "")
		for item, err := range iterjson.Parse(data) {
			errorz.MustZ(err)
			b.Add(item.Key, item.Token)
		}
		out := errorz.Must(b.Bytes())
		fmt.Printf("\n--- minify ---\n%s\n----------\n", out)
	}
	{
		// 🦋Example: reformat json
		b := iterjson.NewBuilder("👉   ", "\t")
		for item, err := range iterjson.Parse(data) {
			errorz.MustZ(err)
			b.Add(item.Key, item.Token)
		}
		out := errorz.Must(b.Bytes())
		fmt.Printf("\n--- reformat ---\n%s\n----------\n", out)
	}
}
