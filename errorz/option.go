package errorz

import (
	"fmt"

	"github.com/ezpkg/stacktracez"
)

type Option struct {
	NoStack     bool
	CallersSkip int
}

func NoStack() Option {
	return Option{NoStack: true}
}
func CallersSkip(n int) Option {
	return Option{CallersSkip: n}
}

func (opt Option) New(msg string) error {
	zErr := &zError{
		msg: msg,
	}
	if !opt.NoStack {
		zErr.stack = stacktracez.StackTraceSkip(opt.CallersSkip + 1)
	}
	return zErr
}

func (opt Option) Newf(format string, args ...any) error {
	zErr := &zError{
		msg: fmt.Sprintf(format, args...),
	}
	if !opt.NoStack {
		zErr.stack = stacktracez.StackTraceSkip(opt.CallersSkip + 1)
	}
	return zErr
}

func (opt Option) Error(msg string) error {
	zErr := &zError{
		msg: msg,
	}
	if !opt.NoStack {
		zErr.stack = stacktracez.StackTraceSkip(opt.CallersSkip + 1)
	}
	return zErr
}

func (opt Option) Errorf(format string, args ...any) error {
	zErr := &zError{
		msg: fmt.Sprintf(format, args...),
	}
	if !opt.NoStack {
		zErr.stack = stacktracez.StackTraceSkip(opt.CallersSkip + 1)
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
		zErr.stack = stacktracez.StackTraceSkip(opt.CallersSkip + 1)
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
		zErr.stack = stacktracez.StackTraceSkip(opt.CallersSkip + 1)
	}
	return zErr
}

func (opt Option) Append(pErr *error, errs ...error) {
	appendErrs(opt, pErr, errs...)
}

func (opt Option) Appendf(pErr *error, err error, msgArgs ...any) {
	err = formatMsg(err, msgArgs)
	if err != nil {
		appendErrs(opt, pErr, err)
	}
}
