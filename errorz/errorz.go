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
		msg:   fmt.Sprintf(format, args...),
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
		msg:   fmt.Sprintf(format, args...),
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
	return fmt.Sprintf("%s", e)
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
