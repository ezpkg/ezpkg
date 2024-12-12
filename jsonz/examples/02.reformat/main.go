package main

import (
	"fmt"
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
		// 🐝Example: minify json
		b := jsonz.NewBuilder("", "")
		for item, err := range jsonz.Parse(data) {
			errorz.MustZ(err)
			b.AddRaw(item.Key, item.Token)
		}
		out := errorz.Must(b.Bytes())
		fmt.Printf("\n--- minify ---\n%s\n----------\n", out)
	}
	{
		// 🦋Example: reformat json
		b := jsonz.NewBuilder("→   ", "\t")
		for item, err := range jsonz.Parse(data) {
			errorz.MustZ(err)
			b.AddRaw(item.Key, item.Token)
		}
		out := errorz.Must(b.Bytes())
		fmt.Printf("\n--- reformat ---\n%s\n----------\n", out)
	}
}
