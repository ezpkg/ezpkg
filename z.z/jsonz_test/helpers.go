package test

import (
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
)

var isDebug = getEnvBool("TESTING_DEBUG")

func getEnvBool(env string) bool {
	str := os.Getenv(env)
	if str != "" {
		return must(strconv.ParseBool(str))
	}
	return false
}

var testcases = []*Testcase{
	// test data from https://github.com/valyala/fastjson
	testcase("small.json"),
	testcase("medium.json"),
	testcase("large.json"),

	// test data from https://github.com/serde-rs/json-benchmark
	testcase("canada.json.gz"),
	testcase("citm_catalog.json.gz"),
	testcase("twitter.json.gz"),

	// test data from https://github.com/Tencent/rapidjson
	testcase("rapid.json.gz"),
}

// Testcase encapsulates the input and states during a test. Testcase is cloned before running each benchmark, and will
// only be accessed locally in each goroutine for parallel benchmarks.
type Testcase struct {
	path string
	data []byte

	// options and states
	reader   *bytes.Reader       // resettable reader
	parallel bool                // read-only
	state    any                 // read-only
	initFn   func(*Testcase) any // called once for each goroutine
}

func testcase(path string) *Testcase {
	var r io.Reader = must(os.Open(filepath.Join("data", path)))
	if strings.HasSuffix(path, ".gz") {
		r = must(gzip.NewReader(r))
	}
	data := must(io.ReadAll(r))
	return &Testcase{
		path: path,
		data: data,
	}
}

func (t *Testcase) State() any        { return t.state }
func (t *Testcase) Reader() io.Reader { return t.reader }
func (t *Testcase) IsParallel() bool  { return t.parallel }

func (t *Testcase) WithParallel(b bool) *Testcase {
	cloned := *t
	cloned.parallel = b
	return &cloned
}

// WithInitFunc clones the Testcase and register an init state function. The function will be run once for each goroutine
// to init the local state for that goroutine. The benchmark code can access the state while running.
func (t *Testcase) WithInitFunc(fn func(tc *Testcase) any) *Testcase {
	cloned := *t
	cloned.initFn = fn
	return &cloned
}

// setupInitState clones the Testcase and runs the registered init function to init the state.
func (t *Testcase) setupInitState() *Testcase {
	cloned := *t
	cloned.reader = bytes.NewReader(t.data)
	if cloned.initFn != nil {
		cloned.state = cloned.initFn(t)
	} else {
		cloned.state = nil
	}
	return &cloned
}

// reset resets the reader
func (t *Testcase) reset() {
	must(t.reader.Seek(0, 0))
}

func benchSkip(b *testing.B, name string, tc *Testcase, fn func(*Testcase) error) error {
	return nil // do nothing
}

// bench runs the benchmark with the given name and testcase. The testcase is cloned for each goroutine.
func bench(b *testing.B, name string, tc *Testcase, fn func(tc *Testcase) error) error {
	parallel := xif(tc.parallel, "parallel/", "")
	name = parallel + name + "/" + tc.path
	var atomErr atomic.Pointer[error]
	b.Run(name, func(b *testing.B) {
		N := xif(tc.parallel, runtime.NumCPU(), 1)
		tcases := make([]*Testcase, N)
		for i := 0; i < N; i++ {
			tcases[i] = tc.setupInitState()
		}
		var atomIdx atomic.Int32
		atomIdx.Store(-1)

		b.ReportAllocs()
		b.SetBytes(int64(len(tc.data)))
		b.ResetTimer()
		if tc.parallel {
			b.RunParallel(func(pb *testing.PB) {
				idx := atomIdx.Add(1)
				tcase := tcases[idx]
				for pb.Next() {
					tcase.reset()
					err := fn(tcase)
					atomErr.CompareAndSwap(nil, &err)
				}
			})
		} else {
			tcase := tcases[0]
			for i := 0; i < b.N; i++ {
				tcase.reset()
				err := fn(tcase)
				atomErr.CompareAndSwap(nil, &err)
			}
		}
		b.StopTimer()
	})
	if err := atomErr.Load(); err != nil {
		return *err
	}
	return nil
}

func shortenError(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()
	if len(msg) <= 100 {
		return msg
	}
	return msg[0:50] + " ... " + msg[len(msg)-50:]
}

func debugf(tb testing.TB, format string, args ...any) {
	if isDebug {
		tb.Logf(format, args...)
	}
}

func xif[T any](cond bool, a, b T) T {
	if cond {
		return a
	} else {
		return b
	}
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
