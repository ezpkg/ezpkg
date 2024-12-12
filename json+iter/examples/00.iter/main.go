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

	// 🎄Example: iterate over json
	for item, err := range jsoniter.Parse(data) {
		errorz.MustZ(err)
		fmt.Printf("%20v : %22v   . %v . %v\n", item.Key, item.Token, item.Level, item.GetPath())
	}
}
