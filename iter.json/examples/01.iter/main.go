package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"ezpkg.io/errorz"
	iterjson "ezpkg.io/iter.json"
)

func main() {
	_, file, _, _ := runtime.Caller(0)
	data := errorz.Must(os.ReadFile(filepath.Dir(file) + "/../alice.json"))

	// 🎄Example: iterate over json
	fmt.Printf("| %25v | %18v | %22v |%v|\n", "PATH", "KEY", "TOKEN", "LVL")
	fmt.Printf("| %s | %s | %s | - |\n", strings.Repeat("-", 25), strings.Repeat("-", 18), strings.Repeat("-", 22))

	for item, err := range iterjson.Parse(data) {
		errorz.MustZ(err)

		fmt.Printf("| %25v | %18v | %22v | %v |\n", item.GetPathString(), item.Key, item.Token, item.Level)
	}
}
