package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"ezpkg.io/errorz"
	jsoniter "ezpkg.io/json+iter"
)

func main() {
	_, file, _, _ := runtime.Caller(0)
	data := errorz.Must(os.ReadFile(filepath.Dir(file) + "/../alice.json"))
	{
		// 🦁Example: filter and output json
		b := jsoniter.NewBuilder("", "    ")
		b.SetSkipEmptyStructures(true) // 👉 skip empty [] or {}
		regexPetName := regexp.MustCompile("pets.*name")
		for item, err := range jsoniter.Parse(data) {
			errorz.MustZ(err)
			if item.Token.IsOpen() || item.Token.IsClose() {
				b.AddRaw(item.Key, item.Token)
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

			b.AddRaw(item.Key, item.Token)
		}
		out := errorz.Must(b.Bytes())
		fmt.Printf("\n--- filter: output json ---\n%s\n----------\n", out)
	}
}
