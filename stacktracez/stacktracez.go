package stacktracez

import (
	"fmt"
	"io"
	"path"
	"runtime"
	"strings"
)

type Frames runtime.Frames
type Frame runtime.Frame

type StackTracer interface {
	StackTrace() *runtime.Frames
}
type StackTracerZ interface {
	StackTraceZ() *Frames
}

func (fz *Frames) StackTraceZ() *Frames        { return fz }
func (fz *Frames) StackTrace() *runtime.Frames { return fz.Unwrap() }

func StackTrace() *Frames {
	var pc [32]uintptr
	runtime.Callers(1, pc[:])
	frames := runtime.CallersFrames(pc[:])
	return (*Frames)(frames)
}

func StackTraceSkip(skip int) *Frames {
	var pc [32]uintptr
	n := runtime.Callers(skip+1, pc[:])
	frames := runtime.CallersFrames(pc[:n])
	return (*Frames)(frames)
}

func Wrap(frames *runtime.Frames) *Frames {
	return (*Frames)(frames)
}

func (fz *Frames) Unwrap() *runtime.Frames {
	return (*runtime.Frames)(fz)
}

func (fz *Frames) Next() (frame runtime.Frame, more bool) {
	return fz.Unwrap().Next()
}

func (fz *Frames) Format(s fmt.State, verb rune) {
	switch verb {
	case 's', 'v':
		switch {
		case s.Flag('#'):
			fprintf(s, "%#v", fz.Unwrap())
		default:
			for fr, ok := fz.Next(); ok; fr, ok = fz.Next() {
				Frame(fr).Format(s, verb)
				writeString(s, "\n")
			}
		}
	}
}

func (f Frame) Format(s fmt.State, verb rune) {
	switch verb {
	case 's', 'v':
		switch {
		case s.Flag('+'):
			fprintf(s, "%s\n\t%s:%d", f.Function, f.File, f.Line)
		case s.Flag('#'):
			fprintf(s, "%#v", f)
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
	}
}

func writeString(w fmt.State, s string) {
	_, _ = io.WriteString(w, s)
}

func fprintf(w fmt.State, format string, args ...any) {
	_, _ = fmt.Fprintf(w, format, args...)
}
