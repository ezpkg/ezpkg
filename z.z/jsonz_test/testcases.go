package test

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Testcase struct {
	Name string
	Data []byte
	Bad  bool // true if it's a bad case
}

var SimpleSet = func() (out []Testcase) {
	list := must(os.ReadDir(filepath.Join(currentDir, "jsonchecker")))
	for _, item := range list {
		name := item.Name()
		if strings.HasSuffix(name, ".json") {
			tcase := load("jsonchecker/" + name)
			tcase.Bad = strings.HasPrefix(name, "fail")
			out = append(out, tcase)
		}
	}
	return out
}()

var LargeSet = []Testcase{
	// test data from https://github.com/valyala/fastjson
	load("data/small.json"),
	load("data/medium.json"),
	load("data/large.json"),

	// test data from https://github.com/serde-rs/json-benchmark
	load("data/canada.json.gz"),
	load("data/citm_catalog.json.gz"),
	load("data/twitter.json.gz"),

	// test data from https://github.com/Tencent/rapidjson
	load("data/rapid.json.gz"),
}

func load(path string) Testcase {
	var r io.Reader = must(os.Open(filepath.Join(currentDir, path)))
	if strings.HasSuffix(path, ".gz") {
		r = must(gzip.NewReader(r))
	}
	data := must(io.ReadAll(r))
	return Testcase{Name: filepath.Base(path), Data: data}
}
