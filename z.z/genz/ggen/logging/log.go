package logging

import (
	"context"
	"runtime"
	"strconv"
	"time"
)

type Level int

const (
	DebugLevel Level = -4
	InfoLevel  Level = 0
	WarnLevel  Level = 4
	ErrorLevel Level = 8
)

func (l Level) String() string {
	switch {
	case l < DebugLevel:
		return "DEBUG-" + strconv.Itoa(int(DebugLevel-l))
	case l == DebugLevel:
		return "DEBUG"
	case l < InfoLevel:
		return "DEBUG+" + strconv.Itoa(int(l-DebugLevel))
	case l == InfoLevel:
		return "INFO"
	case l < WarnLevel:
		return "INFO+" + strconv.Itoa(int(l-InfoLevel))
	case l == WarnLevel:
		return "WARN"
	case l < ErrorLevel:
		return "WARN+" + strconv.Itoa(int(l-WarnLevel))
	case l == ErrorLevel:
		return "ERROR"
	default:
		return "ERROR+" + strconv.Itoa(int(l-ErrorLevel))
	}
}

type Attr struct {
	Key   string
	Value any
}

type Record struct {
	Time    time.Time
	Message string
	Level   Level
	Context context.Context

	pc    uintptr
	attrs []Attr
}

func NewRecord(t time.Time, level Level, msg string, ctx context.Context, attrs []Attr) Record {
	return Record{
		Time:    t,
		Message: msg,
		Level:   level,
		Context: ctx,
		pc:      pc(3),
		attrs:   attrs,
	}
}

func (r Record) Attrs(fn func(Attr)) {
	for _, attr := range r.attrs {
		fn(attr)
	}
}

type Handler interface {
	Enabled(Level) bool
	Handle(Record) error
	WithAttrs([]Attr) Handler
}

type Logger interface {
	Debug(msg string, args ...any)
	Enabled(level Level) bool
	Error(msg string, err error, args ...any)
	Info(msg string, args ...any)
	Log(level Level, msg string, args ...any)
	Warn(msg string, args ...any)
	With(args ...any) Logger
	WithContext(ctx context.Context) Logger
}

type defaultLogger struct {
	handler Handler
	ctx     context.Context
}

func NewLogger(handler Handler) Logger {
	return &defaultLogger{
		handler: handler,
	}
}

func (l defaultLogger) Debug(msg string, args ...any) {
	l.Log(DebugLevel, msg, args...)
}

func (l defaultLogger) Enabled(level Level) bool {
	return l.handler.Enabled(level)
}

func (l defaultLogger) Warn(msg string, args ...any) {
	l.Log(WarnLevel, msg, args...)
}

func (l defaultLogger) Error(msg string, err error, args ...any) {
	if err != nil {
		args = append(args, Attr{Key: "err", Value: err})
	}
	l.Log(ErrorLevel, msg, args...)
}

func (l defaultLogger) Info(msg string, args ...any) {
	l.Log(InfoLevel, msg, args...)
}

func (l defaultLogger) Log(level Level, msg string, args ...any) {
	if !l.Enabled(level) {
		return
	}
	record := NewRecord(time.Now(), level, msg, l.ctx, argsToAttrs(args))
	_ = l.handler.Handle(record)
}

func (l defaultLogger) With(args ...any) Logger {
	attrs := argsToAttrs(args)
	return NewLogger(l.handler.WithAttrs(attrs))
}

func (l defaultLogger) WithContext(ctx context.Context) Logger {
	l.ctx = ctx
	return l
}

// pc returns the program counter at the given stack depth.
func pc(depth int) uintptr {
	var pcs [1]uintptr
	runtime.Callers(depth, pcs[:])
	return pcs[0]
}

func argsToAttrs(args []any) (attrs []Attr) {
	for len(args) > 0 {
		switch arg := args[0].(type) {
		case string:
			if len(args) == 1 {
				attrs = append(attrs, Attr{Key: "BADKEY", Value: arg})
				args = args[1:]
			} else {
				attrs = append(attrs, Attr{Key: arg, Value: args[1]})
				args = args[2:]
			}
		case Attr:
			attrs = append(attrs, arg)
			args = args[1:]
		default:
			attrs = append(attrs, Attr{Key: "BADKEY", Value: arg})
			args = args[1:]
		}
	}
	return attrs
}
