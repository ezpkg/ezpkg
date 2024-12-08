package test

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var mapTestcases = map[string]Testcase{}

type Testcase struct {
	Name string
	Data []byte
	Bad  bool // true if it's a bad case

	ExpectTokens string
	ExpectParse  string
	ExpectFormat string
}

var SimpleSet = func() (out []Testcase) {
	re := regexp.MustCompile(`^[^.]+\.json$`)
	list := must(os.ReadDir(filepath.Join(currentDir, "jsonchecker")))
	for _, item := range list {
		name := item.Name()
		if re.MatchString(name) {
			tcase := load("jsonchecker/" + name)
			tcase.Bad = strings.Contains(name, "fail")
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

func GetTestcase(name string) Testcase {
	tcase, ok := mapTestcases[name]
	if !ok {
		panicf("unknown testcase: %v", name)
	}
	return tcase
}

func load(path string) Testcase {
	loadExt := func(path string, ext string) string {
		path = strings.Replace(path, ".json", ext, 1)
		if _, err := os.Stat(filepath.Join(currentDir, path)); err != nil {
			return ""
		}
		data := must(os.ReadFile(filepath.Join(currentDir, path)))
		return strings.TrimSpace(string(data))
	}

	var r io.Reader = must(os.Open(filepath.Join(currentDir, path)))
	if strings.HasSuffix(path, ".gz") {
		r = must(gzip.NewReader(r))
	}
	data := must(io.ReadAll(r))
	name := filepath.Base(path)
	tcase := Testcase{
		Name: name, Data: data,
		ExpectTokens: loadExt(path, ".token"),
		ExpectParse:  loadExt(path, ".parse"),
		ExpectFormat: loadExt(path, ".format"),
	}
	mapTestcases[name] = tcase
	return tcase
}
