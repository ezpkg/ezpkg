package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"ezpkg.io/errorz"
	iterjson "ezpkg.io/iter.json"
)

func main() {
	_, file, _, _ := runtime.Caller(0)
	data := errorz.Must(os.ReadFile(filepath.Dir(file) + "/../order.json"))
	{
		// 🥝Example: convert all number ids to string
		b := iterjson.NewBuilder("", "    ")
		for item, err := range iterjson.Parse(data) {
			errorz.MustZ(err)
			key, _ := item.GetRawPath().Last().ObjectKey()
			if strings.HasSuffix(key, "_id") {
				id, err0 := item.Token.GetInt()
				if err0 == nil {
					b.Add(item.Key, strconv.Itoa(id))
					continue
				}
			}
			b.Add(item.Key, item.Token)
		}
		out := errorz.Must(b.Bytes())
		fmt.Printf("\n--- convert number id ---\n%s\n----------\n", out)
	}
}
