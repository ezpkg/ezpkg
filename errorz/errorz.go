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
	return (Option{Skip: 1}).Wrap(err, msg)
}

func Wrapf(err error, format string, args ...any) error {
	return (Option{Skip: 1}).Wrapf(err, format, args...)
}

type Option struct {
	NoStack bool
	Skip    int
}

func NoStack() Option {
	return Option{NoStack: true}
}
func CallersSkip(n int) Option {
	return Option{Skip: n}
}

func (opt Option) New(msg string) error {
	zErr := &zError{
		msg: msg,
	}
	if !opt.NoStack {
		zErr.stack = stacktracez.StackTraceSkip(opt.Skip + 1)
	}
	return zErr
}

func (opt Option) Newf(format string, args ...any) error {
	zErr := &zError{
		msg: fmt.Sprintf(format, args...),
	}
	if !opt.NoStack {
		zErr.stack = stacktracez.StackTraceSkip(opt.Skip + 1)
	}
	return zErr
}

func (opt Option) Error(msg string) error {
	zErr := &zError{
		msg: msg,
	}
	if !opt.NoStack {
		zErr.stack = stacktracez.StackTraceSkip(opt.Skip + 1)
	}
	return zErr
}

func (opt Option) Errorf(format string, args ...any) error {
	zErr := &zError{
		msg: fmt.Sprintf(format, args...),
	}
	if !opt.NoStack {
		zErr.stack = stacktracez.StackTraceSkip(opt.Skip + 1)
	}
	return zErr
}

func (opt Option) Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	zErr := &zError{
		msg:   msg,
		cause: err,
	}
	if opt.NoStack {
		return zErr
	}
	stack, ok := err.(stacktracez.StackTracerZ)
	if ok && stack.StackTraceZ() != nil {
		zErr.stack = stack.StackTraceZ()
	} else {
		zErr.stack = stacktracez.StackTraceSkip(opt.Skip + 1)
	}
	return zErr
}

func (opt Option) Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	zErr := &zError{
		msg:   fmt.Sprintf(format, args...),
		cause: err,
	}
	if opt.NoStack {
		return zErr
	}
	stack, ok := err.(stacktracez.StackTracerZ)
	if ok && stack.StackTraceZ() != nil {
		zErr.stack = stack.StackTraceZ()
	} else {
		zErr.stack = stacktracez.StackTraceSkip(opt.Skip + 1)
	}
	return zErr
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
	case 'v':
		e.writeMessage(s)
		if (s.Flag('+') || s.Flag('#')) && e.stack != nil {
			writeString(s, "\n")
			e.stack.Format(s, v)
		}
	case 's':
		e.writeMessage(s)
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
	return e.stack.StackTraceZ()
}

func writeString(w fmt.State, s string) {
	_, _ = io.WriteString(w, s)
}

func fprintf(w fmt.State, format string, args ...any) {
	_, _ = fmt.Fprintf(w, format, args...)
}
