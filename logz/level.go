package logz

import (
	"context"
	"log/slog"
)

var globalLevel = &LevelVar{}

const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

type Level = slog.Level
type LevelVar = slog.LevelVar
type Leveler = slog.Leveler

type Enabler interface {
	Enabled(level Level) bool
}
type CtxEnabler interface {
	Enabled(ctx context.Context, level Level) bool
}

type EnablerFunc func(level Level) bool
type CtxEnablerFunc func(ctx context.Context, level Level) bool

func (fn EnablerFunc) Enabled(level Level) bool {
	return fn(level)
}
func (fn CtxEnablerFunc) Enabled(ctx context.Context, level Level) bool {
	return fn(ctx, level)
}

type Option struct {
	enabler func(ctx context.Context, level Level) bool
}

func (o Option) WithEnabler(enabler EnablerFunc) Option {
	o.enabler = func(ctx context.Context, level Level) bool {
		return enabler(level)
	}
	return o // cloned
}
func (o Option) WithCtxEnabler(enabler CtxEnablerFunc) Option {
	o.enabler = enabler
	return o // cloned
}
func (o Option) WithLevel(level Level) Option {
	o.enabler = func(_ context.Context, l Level) bool { return l >= level }
	return o // cloned
}
func (o Option) WithLeveler(level Leveler) Option {
	o.enabler = func(_ context.Context, l Level) bool { return l >= level.Level() }
	return o // cloned
}
func (o Option) FromLoggerP(logger LoggerP) Logger {
	return &pLogger{l: logger, fn: enabler(o.enabler, logger)}
}
func (o Option) FromLoggerI(logger LoggerI) Logger {
	return &xLogger{w: wrapW{logger}, fn: enabler(o.enabler, logger)}
}
func (o Option) FromLoggerw(logger Loggerw) Logger {
	switch logger := logger.(type) {
	case Logger:
		return logger
	case Loggerx:
		return &xLogger{w: logger, f: logger, fn: enabler(o.enabler, logger)}
	default:
		return &xLogger{w: logger, fn: enabler(o.enabler, logger)}
	}
}
func (o Option) FromLoggerf(logger Loggerf) Logger {
	switch logger := logger.(type) {
	case Logger:
		return logger
	case Loggerx:
		return &xLogger{w: logger, f: logger, fn: enabler(o.enabler, logger)}
	default:
		return &xLogger{f: logger, fn: enabler(o.enabler, logger)}
	}
}
func (o Option) FromLoggerx(logger Loggerx) Logger {
	switch logger := logger.(type) {
	case Logger:
		return logger
	default:
		return &xLogger{w: logger, f: logger, fn: enabler(o.enabler, logger)}
	}
}

func WithEnabler(enabler EnablerFunc) Option {
	return Option{}.WithEnabler(enabler)
}
func WithCtxEnabler(enablerFunc CtxEnablerFunc) Option {
	return Option{}.WithCtxEnabler(enablerFunc)
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
	return globalLevel.Level()
}

func enabler(fn CtxEnablerFunc, logger any) CtxEnablerFunc {
	if fn != nil {
		return fn
	}
	if logger != nil {
		switch x := (logger).(type) {
		case Enabler:
			return func(ctx context.Context, level Level) bool {
				return x.Enabled(level)
			}
		case CtxEnabler:
			return x.Enabled
		}
	}
	return defaultEnabler
}

func defaultEnabler(_ context.Context, level Level) bool {
	return level >= globalLevel.Level()
}
