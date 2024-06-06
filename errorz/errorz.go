package errorz

import (
	"fmt"
	"io"

	"github.com/ezpkg/stacktracez"
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

func (e *zError) Format(s fmt.State, v rune) {
	if e == nil {
		writeString(s, "<nil>")
		return
	}
	switch v {
	case 's', 'v':
		e.writeMessage(s)
		if (s.Flag('+') || s.Flag('#')) && e.stack != nil {
			writeString(s, "\n")
			if e.stack != nil {
				e.stack.Format(s, v)
			}
		}
	case 'q':
		switch {
		case e.msg != "":
			fprintf(s, "%q", e.msg)
		case e.cause != nil:
			fprintf(s, "%q", e.cause)
		default:
			writeString(s, "<empty>")
		}
	}
}

func (e *zError) writeMessage(s fmt.State) {
	switch {
	case e.msg != "" && e.cause != nil:
		fprintf(s, "%s: %v", e.msg, e.cause)
	case e.msg != "":
		writeString(s, e.msg)
	case e.cause != nil:
		fprintf(s, "%v", e.cause)
	default:
		writeString(s, "<empty>")
	}
}

func (e *zError) StackTraceZ() *stacktracez.Frames {
	if e == nil || e.stack == nil {
		return nil
	}
	return e.stack.StackTraceZ()
}

func writeString(w fmt.State, s string) {
	_, _ = io.WriteString(w, s)
}

func fprintf(w fmt.State, format string, args ...any) {
	_, _ = fmt.Fprintf(w, format, args...)
}
