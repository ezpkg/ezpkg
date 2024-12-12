package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"ezpkg.io/errorz"
	"ezpkg.io/jsonz"
)

func main() {
	_, file, _, _ := runtime.Caller(0)
	data := errorz.Must(os.ReadFile(filepath.Dir(file) + "/../alice.json"))
	{
		// 🐳Example: filter and print, use GetPathString()
		fmt.Printf("\n--- filter: print line number ---\n")
		i, rePetName := 0, regexp.MustCompile("pets.*name")
		for item, err := range jsonz.Parse(data) {
			i++
			errorz.MustZ(err)

			path := item.GetPathString()
			switch {
			case path == "name",
				path == "email",
				path == "phone",
				rePetName.MatchString(path),
				strings.Contains(path, "address"):
				// continue
			default:
				continue
			}

			// 👉 print with line number
			fmt.Printf("%2d %20s . %s\n", i, item.Token, item.GetPath())
		}
	}
	{
		// 🐳Example: filter and print, use GetRawPath() and Match()
		fmt.Printf("\n--- filter: print line number ---\n")
		i := 0
		for item, err := range jsonz.Parse(data) {
			i++
			errorz.MustZ(err)

			path := item.GetRawPath()
			switch {
			case path.Match("name"),
				path.Match("email"),
				path.Match("phone"),
				path.Contains("pets") && path.Contains("name"),
				path.Contains("address"):
				// continue
			default:
				continue
			}

			// 👉 print with line number
			fmt.Printf("%2d %20s . %s\n", i, item.Token, item.GetPath())
		}
	}
}
