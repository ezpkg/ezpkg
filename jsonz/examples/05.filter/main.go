package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"ezpkg.io/errorz"
	"ezpkg.io/jsonz"
)

func main() {
	_, file, _, _ := runtime.Caller(0)
	data := errorz.Must(os.ReadFile(filepath.Dir(file) + "/../alice.json"))
	{
		i := 0
		for item, err := range jsonz.Parse(data) {
			i++
			errorz.MustZ(err)

			path := item.GetPathString()
			switch path {
			case "name", "email", "phone":
				// continue
			default:
				if strings.Contains(path, "address") {
					// continue
				} else {
					continue
				}
			}

			// 👉 print with line number
			fmt.Printf("%2d %20s . %s\n", i, item.Token, item.GetPath())
		}
	}
}
