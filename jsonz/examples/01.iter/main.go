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

	for item, err := range jsonz.Parse(data) {
		errorz.MustZ(err)
		fmt.Printf("%20v : %22v   . %v . %v\n", item.Key, item.Token, item.Level, item.GetPath())
	}
}
