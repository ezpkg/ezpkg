package logz

import (
	"context"
	"encoding"
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"runtime"
	"slices"
	"strconv"
	"sync"
	"time"

	"ezpkg.io/bytez"
	"ezpkg.io/colorz"
	"ezpkg.io/stringz"
	"ezpkg.io/typez"
	"ezpkg.io/unsafez"
)

type TextHandlerOptions struct {
	AddSource bool // TODO
	Level     Leveler

	// TODO
	// ReplaceAttr func(groups []string, a Attr) Attr

	FormatTime  func(b []byte, time time.Time) []byte
	FormatLevel func(b []byte, level Level) []byte
}

func (o TextHandlerOptions) ToSlogOptions() slog.HandlerOptions {
	return slog.HandlerOptions{
		AddSource: o.AddSource,
		Level:     o.Level,

		// TODO
		// ReplaceAttr: o.ReplaceAttr,
	}
}

func NewTextHandler(w io.Writer, opt *TextHandlerOptions) Handler {
	if opt == nil {
		opt = &TextHandlerOptions{
			FormatLevel: FormatLevelColor(nil),
		}
	}
	if opt.Level == nil {
		opt.Level = LevelInfo
	}
	h := &textHandler{
		w:  w,
		mu: &sync.Mutex{},

		// replace:  opt.ReplaceAttr,
		addSrc:   opt.AddSource,
		leveler:  opt.Level,
		fmtTime:  opt.FormatTime,
		fmtLevel: opt.FormatLevel,
	}
	return h
}

type textHandler struct {
	w     io.Writer
	mu    *sync.Mutex
	buf   bytez.Buffer
	group []string
	attrs []slog.Attr

	addSrc   bool
	color    bool
	leveler  Leveler
	replace  func(groups []string, a Attr) Attr
	fmtTime  func(b []byte, t time.Time) []byte
	fmtLevel func(b []byte, lvl Level) []byte
}

func (h *textHandler) clone() Handler {
	h2 := *h
	h2.buf = bytez.Buffer{}
	h2.group = slices.Clip(h2.group)
	h2.attrs = slices.Clip(h2.attrs)
	return &h2
}

func (h *textHandler) syncWrite(b []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.w.Write(b)
	return err
}

func (h *textHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.leveler.Level() <= level
}

func (h *textHandler) Handle(ctx context.Context, r slog.Record) error {
	var b bytez.Bytes
	if h.fmtTime != nil {
		b = h.fmtTime(b, r.Time)
		b.WriteByteZ('\t')
	}
	if h.fmtLevel != nil {
		b = h.fmtLevel(b, r.Level)
		b.WriteByteZ('\t')
	}
	if h.addSrc && r.PC != 0 {
		b = appendSrc(b, r.PC)
		b.WriteByteZ('\t')
	}
	b.WriteStringZ(r.Message)

	isFirst := true
	r.Attrs(func(a slog.Attr) bool {
		if isFirst {
			b.WriteByteZ('\t')
			isFirst = false
		} else {
			b.WriteByteZ(' ')
		}
		if a.Key == "" {
			return true
		}
		b.WriteStringZ(a.Key)
		b.WriteStringZ("=")
		b = appendValue(b, a.Value)
		return true
	})
	b.WriteByteZ('\n')
	return h.syncWrite(b)
}

func (h *textHandler) WithAttrs(attrs []slog.Attr) Handler {
	h = h.clone().(*textHandler)
	h.attrs = append(h.attrs, attrs...)
	return h
}

func (h *textHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	h = h.clone().(*textHandler)
	h.group = append(h.group, name)
	return h
}

type FormatLevelColorOptions struct {
	DebugColor string // colorz.Magenta.Code()
	InfoColor  string // colorz.Cyan.Code()
	WarnColor  string // colorz.Yellow.Code()
	ErrorColor string // colorz.Red.Code()
}

func FormatLevelColor(opt *FormatLevelColorOptions) func([]byte, Level) []byte {
	colorDebug := colorz.Magenta.Code()
	colorInfo := colorz.Cyan.Code()
	colorWarn := colorz.Yellow.Code()
	colorError := colorz.Red.Code()
	colorReset := colorz.Reset.Code()
	if opt != nil {
		colorDebug = typez.Coalesce(opt.DebugColor, colorDebug)
		colorInfo = typez.Coalesce(opt.InfoColor, colorInfo)
		colorWarn = typez.Coalesce(opt.WarnColor, colorWarn)
		colorError = typez.Coalesce(opt.ErrorColor, colorError)
	}
	return func(b []byte, level Level) []byte {
		switch {
		case level <= LevelDebug:
			b = append(b, colorDebug...)
			b = append(b, level.String()...)
			b = append(b, colorReset...)
		case level <= LevelInfo:
			b = append(b, colorInfo...)
			b = append(b, level.String()...)
			b = append(b, colorReset...)
		case level <= LevelWarn:
			b = append(b, colorWarn...)
			b = append(b, level.String()...)
			b = append(b, colorReset...)
		default:
			b = append(b, colorError...)
			b = append(b, level.String()...)
			b = append(b, colorReset...)
		}
		return b
	}
}

func appendSrc(b []byte, pc uintptr) []byte {
	src := pcToSource(pc)
	file := stringz.LastNParts(src.File, 2, "/")
	file = typez.Coalesce(file, "<unknown>")
	b = append(b, file...)
	b = append(b, ':')
	b = strconv.AppendInt(b, int64(src.Line), 10)
	return b
}

func appendValue(b []byte, v slog.Value) []byte {
	switch v.Kind() {
	case slog.KindString:
		return appendStr(b, v.String())
	case slog.KindInt64:
		return strconv.AppendInt(b, v.Int64(), 10)
	case slog.KindUint64:
		return strconv.AppendUint(b, v.Uint64(), 10)
	case slog.KindFloat64:
		return strconv.AppendFloat(b, v.Float64(), 'g', -1, 64)
	case slog.KindBool:
		return strconv.AppendBool(b, v.Bool())
	case slog.KindDuration:
		return append(b, v.Duration().String()...)
	case slog.KindTime:
		return append(b, v.Time().String()...)
	case slog.KindGroup:
		return fmt.Append(b, v.Group())

	case slog.KindAny:
		if tm, ok := v.Any().(encoding.TextMarshaler); ok {
			data, err := tm.MarshalText()
			if err != nil {
				return appendError(b, err)
			}
			return appendStr(b, unsafez.BytesToString(data))
		}
		if data, ok := isBytes(v.Any()); ok {
			return appendStr(b, unsafez.BytesToString(data))
		}
		return appendStr(b, fmt.Sprintf("%+v", v.Any()))

	case slog.KindLogValuer:
		panic("should never happen")
	}
	return b
}

func appendStr(b []byte, s string) []byte {
	if needsQuoting(s) {
		b = strconv.AppendQuote(b, s)
	} else {
		b = append(b, s...)
	}
	return b
}

func appendError(b []byte, err error) []byte {
	return append(b, fmt.Sprintf("!ERROR:%v", err)...)
}

func isBytes(a any) ([]byte, bool) {
	if s, ok := a.([]byte); ok {
		return s, true
	}
	t := reflect.TypeOf(a)
	if t != nil && t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Uint8 {
		return reflect.ValueOf(a).Bytes(), true
	}
	return nil, false
}

func pcToSource(pc uintptr) slog.Source {
	fs := runtime.CallersFrames([]uintptr{pc})
	f, _ := fs.Next()
	return slog.Source{
		Function: f.Function,
		File:     f.File,
		Line:     f.Line,
	}
}
