package stacktracez

import (
	"runtime"
	"testing"
)

func BenchmarkStd(b *testing.B) {
	b.ReportAllocs()
	var frames *runtime.Frames
	for i := 0; i < b.N; i++ {
		var pc [32]uintptr
		runtime.Callers(2, pc[:])
		frames = runtime.CallersFrames(pc[:])
	}
	b.StopTimer()
	frames.Next()
}

func BenchmarkStdNext(b *testing.B) {
	b.ReportAllocs()
	var frames *runtime.Frames
	for i := 0; i < b.N; i++ {
		var pc [32]uintptr
		runtime.Callers(2, pc[:])
		frames = runtime.CallersFrames(pc[:])
		_, more := frames.Next()
		for more {
			_, more = frames.Next()
		}
	}
}

func BenchmarkZ(b *testing.B) {
	b.ReportAllocs()
	var frames *Frames
	for i := 0; i < b.N; i++ {
		frames = StackTrace()
	}
	b.StopTimer()
	frames.GetFrames()
}

func BenchmarkZGetFrames(b *testing.B) {
	b.ReportAllocs()
	var frames *Frames
	for i := 0; i < b.N; i++ {
		frames = StackTrace()
		frames.GetFrames()
	}
}
