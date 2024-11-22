package test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync/atomic"
	"testing"
)

var isDebug = getEnvBool("TESTING_DEBUG")
var currentDir = getCurrentDir()

func getEnvBool(env string) bool {
	str := os.Getenv(env)
	if str != "" {
		return must(strconv.ParseBool(str))
	}
	return false
}

// xTestcase encapsulates the input and states during a test. xTestcase is cloned before running each benchmark, and will
// only be accessed locally in each goroutine for parallel benchmarks.
type xTestcase struct {
	Testcase

	// options and states
	reader   *bytes.Reader        // resettable reader
	parallel bool                 // read-only
	state    any                  // read-only
	initFn   func(*xTestcase) any // called once for each goroutine
}

func (t *xTestcase) State() any        { return t.state }
func (t *xTestcase) Reader() io.Reader { return t.reader }
func (t *xTestcase) IsParallel() bool  { return t.parallel }

func (t *xTestcase) WithParallel(b bool) *xTestcase {
	cloned := *t
	cloned.parallel = b
	return &cloned
}

// WithInitFunc clones the Testcase and register an init state function. The function will be run once for each goroutine
// to init the local state for that goroutine. The benchmark code can access the state while running.
func (t *xTestcase) WithInitFunc(fn func(tc *xTestcase) any) *xTestcase {
	cloned := *t
	cloned.initFn = fn
	return &cloned
}

// setupInitState clones the Testcase and runs the registered init function to init the state.
func (t *xTestcase) setupInitState() *xTestcase {
	cloned := *t
	cloned.reader = bytes.NewReader(t.Data)
	if cloned.initFn != nil {
		cloned.state = cloned.initFn(t)
	} else {
		cloned.state = nil
	}
	return &cloned
}

// reset resets the reader
func (t *xTestcase) reset() {
	must(t.reader.Seek(0, 0))
}

func benchSkip(b *testing.B, name string, tc *xTestcase, fn func(*xTestcase) error) error {
	return nil // do nothing
}

// bench runs the benchmark with the given name and testcase. The testcase is cloned for each goroutine.
func bench(b *testing.B, name string, tc *xTestcase, fn func(tc *xTestcase) error) error {
	parallel := xif(tc.parallel, "parallel/", "")
	name = parallel + name + "/" + tc.Name
	var atomErr atomic.Pointer[error]
	b.Run(name, func(b *testing.B) {
		N := xif(tc.parallel, runtime.NumCPU(), 1)
		tcases := make([]*xTestcase, N)
		for i := 0; i < N; i++ {
			tcases[i] = tc.setupInitState()
		}
		var atomIdx atomic.Int32
		atomIdx.Store(-1)

		b.ReportAllocs()
		b.SetBytes(int64(len(tc.Data)))
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

func getCurrentDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("unexpected")
	}
	return filepath.Dir(file)
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
