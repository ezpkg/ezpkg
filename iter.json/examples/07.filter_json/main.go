package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"ezpkg.io/errorz"
	iterjson "ezpkg.io/iter.json"
)

func main() {
	_, file, _, _ := runtime.Caller(0)
	data := errorz.Must(os.ReadFile(filepath.Dir(file) + "/../alice.json"))
	{
		// 🦁Example: filter and output json
		b := iterjson.NewBuilder("", "    ")
		b.SetSkipEmptyStructures(true) // 👉 skip empty [] or {}
		regexPetName := regexp.MustCompile("pets.*name")
		for item, err := range iterjson.Parse(data) {
			errorz.MustZ(err)
			if item.Token.IsOpen() || item.Token.IsClose() {
				b.Add(item.Key, item.Token)
				continue
			}

			path := item.GetPathString()
			switch {
			case path == "name",
				path == "email",
				path == "phone",
				regexPetName.MatchString(path),
				strings.Contains(path, "address"):
				// continue
			default:
				continue
			}

			b.Add(item.Key, item.Token)
		}
		out := errorz.Must(b.Bytes())
		fmt.Printf("\n--- filter: output json ---\n%s\n----------\n", out)
	}
}
