package logz

import (
	"ezpkg.io/stringz"
)

type strLevel string

const (
	strDebug strLevel = "DEBUG"
	strInfo  strLevel = " INFO"
	strWarn  strLevel = " WARN"
	strError strLevel = "ERROR"
)

type wrapW struct{ LoggerI }

func (w wrapW) Debugw(msg string, keyValues ...any) { w.LoggerI.Debug(msg, keyValues...) }
func (w wrapW) Infow(msg string, keyValues ...any)  { w.LoggerI.Info(msg, keyValues...) }
func (w wrapW) Warnw(msg string, keyValues ...any)  { w.LoggerI.Warn(msg, keyValues...) }
func (w wrapW) Errorw(msg string, keyValues ...any) { w.LoggerI.Error(msg, keyValues...) }

type zKV struct {
	key string
	val any
}

func formatWf(msg string, level strLevel, kv []any) string {
	var b stringz.Builder
	if level != "" {
		b.WriteStringZ(string(level))
		b.WriteStringZ(": ")
	}
	b.WriteStringZ(msg)
	for i, N := 0, len(kv); i < N; i++ {
		if key, ok := kv[i].(string); ok && i < N-1 {
			b.Printf(" %v=%q", key, kv[i+1])
			i++
		} else {
			b.Printf(" [%d]=%q", i, kv[i])
		}
	}
	if level != "" {
		b.Println()
	}
	return b.String()
}

func formatf(format string, level strLevel, args []any) string {
	var b stringz.Builder
	if level != "" {
		b.WriteStringZ(string(level))
		b.WriteStringZ(": ")
	}
	b.Printf(format, args...)
	if level != "" {
		b.Println()
	}
	return b.String()
}
