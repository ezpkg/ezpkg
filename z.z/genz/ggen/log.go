package ggen

import (
	"fmt"
	"io"

	log "ezpkg.io/genz/ggen/logging"
)

type Logger = log.Logger
type LogAttr = log.Attr
type LogLevel = log.Level
type LogHandler = log.Handler

const (
	DebugLevel = log.DebugLevel
	InfoLevel  = log.InfoLevel
	WarnLevel  = log.WarnLevel
	ErrorLevel = log.ErrorLevel
)

var logger Logger

type embededLogger struct {
	Logger
}

type defaultLogHandler struct {
	w     io.Writer
	level log.Level
	attrs []log.Attr
}

func (h defaultLogHandler) Enabled(level log.Level) bool {
	return level >= h.level
}

func (h defaultLogHandler) WithAttrs(attrs []log.Attr) log.Handler {
	newAttrs := make([]log.Attr, 0, len(h.attrs)+len(attrs))
	newAttrs = append(newAttrs, h.attrs...)
	newAttrs = append(newAttrs, attrs...)
	h.attrs = newAttrs
	return h // copy
}

func (h defaultLogHandler) Handle(r log.Record) (err error) {
	debugEnabled := h.Enabled(-1)
	if debugEnabled {
		_, err = fmt.Fprintf(h.w, "%7s: %s", r.Level, r.Message)
	} else if r.Level != 0 {
		_, err = fmt.Fprintf(h.w, "%s", r.Message)
	} else {
		_, err = fmt.Fprint(h.w, r.Message)
	}
	if err != nil {
		return err
	}
	for _, attr := range h.attrs {
		fmt.Fprintf(h.w, " %s=%v", attr.Key, attr.Value)
	}
	r.Attrs(func(attr log.Attr) {
		fmt.Fprintf(h.w, " %s=%q", attr.Key, attr.Value)
	})
	_, err = fmt.Fprintf(h.w, "\n")
	return err
}
