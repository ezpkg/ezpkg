package errorz // import "ezpkg.io/errorz"

import (
	"fmt"

	"ezpkg.io/fmtz"
	"ezpkg.io/stacktracez"
)

func Must[T any](v T, err error) T {
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
	return v
}
func MustZ(err error) {
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
}
func Must2[A, B any](a A, b B, err error) (A, B) {
	MustZ(err)
	return a, b
}
func Must3[A, B, C any](a A, b B, c C, err error) (A, B, C) {
	MustZ(err)
	return a, b, c
}
func Must4[A, B, C, D any](a A, b B, c C, d D, err error) (A, B, C, D) {
	MustZ(err)
	return a, b, c, d
}
func Must5[A, B, C, D, E any](a A, b B, c C, d D, e E, err error) (A, B, C, D, E) {
	MustZ(err)
	return a, b, c, d, e
}
func Must6[A, B, C, D, E, F any](a A, b B, c C, d D, e E, f F, err error) (A, B, C, D, E, F) {
	MustZ(err)
	return a, b, c, d, e, f
}
func Must7[A, B, C, D, E, F, G any](a A, b B, c C, d D, e E, f F, g G, err error) (A, B, C, D, E, F, G) {
	MustZ(err)
	return a, b, c, d, e, f, g
}
func Must8[A, B, C, D, E, F, G, H any](a A, b B, c C, d D, e E, f F, g G, h H, err error) (A, B, C, D, E, F, G, H) {
	MustZ(err)
	return a, b, c, d, e, f, g, h
}

var _ stacktracez.StackTracerZ = (*zError)(nil)

type zError struct {
	msg   string
	cause error
	stack *stacktracez.Frames
}

func New(msg string) error {
	return &zError{
		msg:   msg,
		stack: stacktracez.StackTraceSkip(1),
	}
}

func Newf(format string, args ...any) error {
	return &zError{
		msg:   sprintf(format, args...),
		stack: stacktracez.StackTraceSkip(1),
	}
}

func Error(msg string) error {
	return &zError{
		msg:   msg,
		stack: stacktracez.StackTraceSkip(1),
	}
}

func Errorf(format string, args ...any) error {
	return &zError{
		msg:   sprintf(format, args...),
		stack: stacktracez.StackTraceSkip(1),
	}
}

// MapFunc maps a value and an error to another value and error using the input function.
func MapFunc[T, Q any](fn func(T) Q) func(v T, err error) (Q, error) {
	return func(v T, err error) (Q, error) {
		if err == nil {
			return fn(v), nil
		}
		var zero Q
		return zero, err
	}
}

func Wrap(err error, msg string) error {
	return (Option{CallersSkip: 1}).Wrap(err, msg)
}

func Wrapf(err error, format string, args ...any) error {
	return (Option{CallersSkip: 1}).Wrapf(err, format, args...)
}

func (e *zError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.cause == nil && e.msg != "" {
		return e.msg
	}
	return sprintf("%s", e)
}

func (e *zError) Format(s0 fmt.State, verb rune) {
	s := fmtz.WrapState(s0)
	if e == nil {
		s.WriteStringZ("<nil>")
		return
	}
	switch verb {
	case 's', 'v':
		e.writeMessage(s, verb)
		if (s.Flag('+') || s.Flag('#')) && e.stack != nil {
			s.WriteStringZ("\n")
			if e.stack != nil {
				e.stack.Format(s, verb)
			}
		}
	case 'q':
		switch {
		case e.msg != "":
			s.Printf("%q", e.msg)
		case e.cause != nil:
			e.formatCause(s, verb)
		default:
			s.WriteStringZ("<empty>")
		}
	}
}

func (e *zError) writeMessage(s fmtz.State, verb rune) {
	switch {
	case e.msg != "" && e.cause != nil:
		s.WriteStringZ(e.msg)
		s.WriteStringZ(": ")
		e.formatCause(s, verb)
	case e.cause != nil:
		e.formatCause(s, verb)
	case e.msg != "":
		s.WriteStringZ(e.msg)
	default:
		s.WriteStringZ("<empty>")
	}
}

func (e *zError) formatCause(s fmtz.State, verb rune) {
	if _, ok := e.cause.(stacktracez.StackTracerZ); ok {
		s.Format(verb, e.cause)
	} else if s.Flag('+') {
		s.Printf("%+v", e.cause)
	} else {
		s.Printf("%v", e.cause)
	}
}

func (e *zError) StackTraceZ() *stacktracez.Frames {
	if e == nil || e.stack == nil {
		return nil
	}
	return e.stack.StackTraceZ()
}
