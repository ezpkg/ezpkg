package errorz // import "ezpkg.io/errorz"

import (
	"fmt"

	"ezpkg.io/fmtz"
	"ezpkg.io/stacktracez"
)

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
func MustZ(err error) {
	if err != nil {
		panic(err)
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

func (e *zError) Format(s0 fmt.State, v rune) {
	s := fmtz.WrapState(s0)
	if e == nil {
		s.WriteStringZ("<nil>")
		return
	}
	switch v {
	case 's', 'v':
		e.writeMessage(s)
		if (s.Flag('+') || s.Flag('#')) && e.stack != nil {
			s.WriteStringZ("\n")
			if e.stack != nil {
				e.stack.Format(s, v)
			}
		}
	case 'q':
		switch {
		case e.msg != "":
			s.Printf("%q", e.msg)
		case e.cause != nil:
			s.Printf("%q", e.cause)
		default:
			s.WriteStringZ("<empty>")
		}
	}
}

func (e *zError) writeMessage(s0 fmt.State) {
	s := fmtz.WrapState(s0)
	switch {
	case e.msg != "" && e.cause != nil:
		s.Printf("%s: %v", e.msg, e.cause)
	case e.msg != "":
		s.WriteStringZ(e.msg)
	case e.cause != nil:
		s.Printf("%v", e.cause)
	default:
		s.WriteStringZ("<empty>")
	}
}

func (e *zError) StackTraceZ() *stacktracez.Frames {
	if e == nil || e.stack == nil {
		return nil
	}
	return e.stack.StackTraceZ()
}
