package stacktracez // import "ezpkg.io/stacktracez"

import (
	"fmt"
	"path"
	"runtime"
	"strings"
	"sync"

	"ezpkg.io/fmtz"
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
	runtime.Callers(2, pc[:])
	frames := runtime.CallersFrames(pc[:])
	return &Frames{frames: frames}
}

func StackTraceSkip(skip int) *Frames {
	var pc [maxFrames]uintptr
	n := runtime.Callers(skip+2, pc[:])
	frames := runtime.CallersFrames(pc[:n])
	return &Frames{frames: frames}
}

func (fz *Frames) Format(s0 fmt.State, verb rune) {
	s := fmtz.WrapState(s0)
	if fz == nil {
		s.WriteStringZ("<nil>")
		return
	}
	switch verb {
	case 's', 'v':
		formatFrames(s, verb, fz.GetFrames())
	}
}

func (fz *Frames) GetFrames() []Frame {
	if fz == nil {
		return nil
	}
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

func (f Frame) Format(s0 fmt.State, verb rune) {
	s := fmtz.WrapState(s0)
	switch verb {
	case 's', 'v':
		switch {
		case s.Flag('+'):
			s.Printf("%s\n\t%s:%d", f.Function, f.File, f.Line)
		default:
			pkg, file, line, fn := f.Components()
			s.Printf("%s/%s:%d · %s", pkg, file, line, fn)
		}
	case 'd':
		s.Printf("%d", f.Line)
	}
}

func (f Frame) Components() (pkg, file string, line int, fn string) {
	sepIdx := strings.LastIndexByte(f.Function, '/')
	if sepIdx < 0 {
		sepIdx = 0
	}
	dotIdx := strings.IndexByte(f.Function[sepIdx:], '.')
	if dotIdx <= 0 {
		dotIdx = -1
	}
	pkg, fn = f.Function[:sepIdx+dotIdx], f.Function[sepIdx+dotIdx+1:]
	return pkg, path.Base(f.File), f.Line, fn
}

func formatFrames(s fmtz.State, verb rune, frames []Frame) {
	if len(frames) == 0 {
		s.WriteStringZ("[]")
	}
	for _, frame := range frames {
		frame.Format(s, verb)
		s.WriteStringZ("\n")
	}
}
