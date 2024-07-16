package logz

import (
	"log/slog"
	"sync/atomic"
)

var globalLevel = &LevelVar{}

const (
	LevelDebug = Level(slog.LevelDebug)
	LevelInfo  = Level(slog.LevelInfo)
	LevelWarn  = Level(slog.LevelWarn)
	LevelError = Level(slog.LevelError)
)

type Level slog.Level

func (l Level) Unwrap() slog.Level { return slog.Level(l) }
func (l Level) ToInt8() int8       { return int8(l / 4) }
func (l Level) ToInt() int         { return int(l) }
func (l Level) String() string     { return l.Unwrap().String() }

type LevelVar struct {
	level atomic.Int64
}

func (l *LevelVar) Set(level Level) {
	l.level.Store(int64(level))
}
func (l *LevelVar) Get() Level {
	return Level(l.level.Load())
}
func (l *LevelVar) Level() Level {
	return Level(l.level.Load())
}
func (l *LevelVar) String() string {
	return l.Get().String()
}

type Leveler interface{ Level() Level }
type EnablerFunc func(level Level) bool

type Option struct {
	enabler func(level Level) bool
}

func (o Option) WithEnabler(enabler EnablerFunc) Option {
	o.enabler = enabler
	return o // cloned
}
func (o Option) WithLevel(level Level) Option {
	o.enabler = func(l Level) bool { return l >= level }
	return o // cloned
}
func (o Option) WithLeveler(level Leveler) Option {
	o.enabler = func(l Level) bool { return l >= level.Level() }
	return o // cloned
}
func (o Option) FromLoggerP(logger LoggerP) Logger {
	return &pLogger{l: logger, fn: enabler(o.enabler)}
}
func (o Option) FromLoggerI(logger LoggerI) Logger {
	return &xLogger{w: wrapW{logger}, fn: enabler(o.enabler)}
}
func (o Option) FromLoggerw(logger Loggerw) Logger {
	switch logger := logger.(type) {
	case Logger:
		return logger
	case Loggerx:
		return &xLogger{w: logger, f: logger, fn: enabler(o.enabler)}
	default:
		return &xLogger{w: logger, fn: enabler(o.enabler)}
	}
}
func (o Option) FromLoggerf(logger Loggerf) Logger {
	switch logger := logger.(type) {
	case Logger:
		return logger
	case Loggerx:
		return &xLogger{w: logger, f: logger, fn: enabler(o.enabler)}
	default:
		return &xLogger{f: logger, fn: enabler(o.enabler)}
	}
}
func (o Option) FromLoggerx(logger Loggerx) Logger {
	switch logger := logger.(type) {
	case Logger:
		return logger
	default:
		return &xLogger{w: logger, f: logger, fn: enabler(o.enabler)}
	}
}

func WithEnabler(enabler EnablerFunc) Option {
	return Option{}.WithEnabler(enabler)
}
func WithLevel(level Level) Option {
	return Option{}.WithLevel(level)
}
func WithLeveler(level Leveler) Option {
	return Option{}.WithLeveler(level)
}
func DefaultLevelVar() *LevelVar {
	return globalLevel
}
func SetDefaultEnableLevel(level Level) {
	globalLevel.Set(level)
}
func GetDefaultEnableLevel() Level {
	return globalLevel.Get()
}

func enabler(fn EnablerFunc) EnablerFunc {
	if fn != nil {
		return fn
	}
	return func(l Level) bool { return l >= globalLevel.Level() }
}
