package stacktracez

import (
	"fmt"
	"io"
	"path"
	"runtime"
	"strings"
	"sync"
)

const maxFrames = 32

type Frames struct {
	frames *runtime.Frames
	mutex  sync.RWMutex
	cached []Frame
}
type Frame runtime.Frame

type StackTracerZ interface {
	StackTraceZ() *Frames
}

func (fz *Frames) StackTraceZ() *Frames { return fz }

func StackTrace() *Frames {
	var pc [maxFrames]uintptr
	runtime.Callers(1, pc[:])
	frames := runtime.CallersFrames(pc[:])
	return &Frames{frames: frames}
}

func StackTraceSkip(skip int) *Frames {
	var pc [maxFrames]uintptr
	n := runtime.Callers(skip+1, pc[:])
	frames := runtime.CallersFrames(pc[:n])
	return &Frames{frames: frames}
}

func (fz *Frames) Format(s fmt.State, verb rune) {
	switch verb {
	case 's', 'v':
		formatFrames(s, verb, fz.GetFrames())
	}
}

func (fz *Frames) GetFrames() []Frame {
	fz.mutex.RLock()
	if cached := fz.cached; cached != nil {
		fz.mutex.RUnlock()
		return cached
	}
	fz.mutex.RUnlock()
	fz.mutex.Lock()
	defer fz.mutex.Unlock()
	if cached := fz.cached; cached != nil {
		return cached
	}
	fz.cached = make([]Frame, 0, maxFrames)
	for fr, ok := fz.frames.Next(); ok; fr, ok = fz.frames.Next() {
		fz.cached = append(fz.cached, Frame(fr))
	}
	return fz.cached
}

func (f Frame) Format(s fmt.State, verb rune) {
	switch verb {
	case 's', 'v':
		switch {
		case s.Flag('+'):
			fprintf(s, "%s\n\t%s:%d", f.Function, f.File, f.Line)
		default:
			sepIdx := strings.LastIndexByte(f.Function, '/')
			if sepIdx < 0 {
				sepIdx = 0
			}
			dotIdx := strings.IndexByte(f.Function[sepIdx:], '.')
			if dotIdx <= 0 {
				dotIdx = -1
			}
			pkgPath := f.Function[:sepIdx+dotIdx]
			fprintf(s, "%s/%s:%d · %s", pkgPath, path.Base(f.File), f.Line, f.Function[sepIdx+dotIdx+1:])
		}
	case 'd':
		fprintf(s, "%d", f.Line)
	}
}

func formatFrames(s fmt.State, verb rune, frames []Frame) {
	for _, frame := range frames {
		frame.Format(s, verb)
		writeString(s, "\n")
	}
}

func writeString(w fmt.State, s string) {
	_, _ = io.WriteString(w, s)
}

func fprintf(w fmt.State, format string, args ...any) {
	_, _ = fmt.Fprintf(w, format, args...)
}
