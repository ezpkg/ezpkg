package stacktracez

import (
	"runtime"
	"testing"
)

func BenchmarkStd(b *testing.B) {
	b.ReportAllocs()
	var frames *runtime.Frames
	for i := 0; i < b.N; i++ {
		var callers [32]uintptr
		frames = runtime.CallersFrames(callers[:])
		frames.Next()
	}
}

func BenchmarkZ(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		frames := StackTrace()
		frames.GetFrames()
	}
}
