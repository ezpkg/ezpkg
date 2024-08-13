// Package logz provides common interfaces and utilities for working with other log packages, including [log/slog] and [zap].
//
// [log/slog]: https://pkg.go.dev/log/slog
// [zap]: https://pkg.go.dev/go.uber.org/zap
package logz // import ezpkg.io/logz

import (
	"context"
	"io"
	"log/slog"
)

type Handler = slog.Handler
type Attr = slog.Attr

type Logger interface {
	Debugw(msg string, keyValues ...any)
	Infow(msg string, keyValues ...any)
	Warnw(msg string, keyValues ...any)
	Errorw(msg string, keyValues ...any)

	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)

	Enabled(ctx context.Context, level Level) bool
	With(keyValues ...any) Logger
}

// LoggerP is implemented by stdlib log.Logger.
type LoggerP interface {
	Printf(format string, args ...any)
}

// LoggerI is implemented by stdlib slog.Logger.
type LoggerI interface {
	Debug(msg string, keyValues ...any)
	Info(msg string, keyValues ...any)
	Warn(msg string, keyValues ...any)
	Error(msg string, keyValues ...any)
}

// Loggerw is implemented by zap.SugaredLogger.
type Loggerw interface {
	Debugw(msg string, keyValues ...any)
	Infow(msg string, keyValues ...any)
	Warnw(msg string, keyValues ...any)
	Errorw(msg string, keyValues ...any)
}

type Loggerf interface {
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
}

type Loggerx interface {
	Loggerw
	Loggerf
}

type logger0ctx interface {
	LoggerI
	DebugContext(ctx context.Context, msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
}

func New(h Handler) Logger {
	opt := Option{
		enabler: func(ctx context.Context, level Level) bool {
			return h.Enabled(ctx, level)
		},
	}
	sl := slog.New(h)
	return opt.FromLoggerI(sl)
}

func DefaultLogger(w io.Writer) Logger {
	return New(NewTextHandler(w, nil))
}

func FromLoggerP(logger LoggerP) Logger {
	return Option{}.FromLoggerP(logger)
}
func FromLoggerI(logger LoggerI) Logger {
	return Option{}.FromLoggerI(logger)
}
func FromLoggerw(logger Loggerw) Logger {
	return Option{}.FromLoggerw(logger)
}
func FromLoggerf(logger Loggerf) Logger {
	return Option{}.FromLoggerf(logger)
}
func FromLoggerx(logger Loggerx) Logger {
	return Option{}.FromLoggerx(logger)
}
