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

func appendKV(list []zKV, kv []any) []zKV {
	list = list[:] // return different slice
	for i, N := 0, len(kv); i < N; i++ {
		if key, ok := kv[i].(string); ok && i < N-1 {
			list = append(list, zKV{key, kv[i+1]})
			i++
		} else {
			list = append(list, zKV{"", kv[i]})
		}
	}
	return list
}

func formatWf(level strLevel, msg string, zkv []zKV, kv []any) stringz.StringFunc {
	return func() string {
		var b stringz.Builder
		if level != "" {
			b.WriteStringZ(string(level))
			b.WriteStringZ(": ")
		}
		b.WriteStringZ(msg)

		extraIdx := 0
		for i, N := 0, len(zkv); i < N; i++ {
			k := zkv[i].key
			if k != "" {
				b.Printf(" %v=", zkv[i].key)
				formatVal(&b, zkv[i].val)
			} else {
				b.Printf(" [%d]=", extraIdx)
				formatVal(&b, zkv[i].val)
				extraIdx++
			}
		}
		for i, N := 0, len(kv); i < N; i++ {
			if key, ok := kv[i].(string); ok && i < N-1 {
				b.Printf(" %v=", key)
				formatVal(&b, kv[i+1])
				i++
			} else {
				b.Printf(" [%d]=", extraIdx)
				formatVal(&b, kv[i])
				extraIdx++
			}
		}
		if level != "" {
			b.Println()
		}
		return b.String()
	}
}

func formatf(level strLevel, format string, args []any, zkv []zKV) stringz.StringFunc {
	return func() string {
		var b stringz.Builder
		if level != "" {
			b.WriteStringZ(string(level))
			b.WriteStringZ(": ")
		}
		b.Printf(format, args...)
		for i, N := 0, len(zkv); i < N; i++ {
			b.Printf(" %v=", zkv[i].key)
			formatVal(&b, zkv[i].val)
		}
		if level != "" {
			b.Println()
		}
		return b.String()
	}
}

func formatVal(b *stringz.Builder, val any) {
	switch v := (val).(type) {
	case string, []byte:
		b.Printf("%q", v)
	default:
		b.Printf("%v", v)
	}
}
