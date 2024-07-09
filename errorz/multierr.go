package errorz

import (
	"fmt"
	"reflect"

	"ezpkg.io/fmtz"
	"ezpkg.io/stacktracez"
)

var _ Errors = (*zErrors)(nil)
var _ stacktracez.StackTracerZ = (*zErrors)(nil)

type Errors interface {
	error
	Errors() []error
	Unwrap() []error
}

func Append(err0 error, errs ...error) error {
	appendErrs(Option{}, &err0, errs...)
	return err0
}

func AppendTo(pErr *error, errs ...error) {
	appendErrs(Option{}, pErr, errs...)
}

func appendErrs(opt Option, pErr *error, errs ...error) {
	switch err0 := (*pErr).(type) {
	case nil:
		var zErrs zErrors
		zErrs.Append(errs...)
		*pErr = zErrs.process(opt)

	case *zErrors:
		err0.Append(errs...)
		*pErr = err0.process(opt)

	default:
		var zErrs zErrors
		zErrs.errors = make([]error, len(errs))
		zErrs.Append(errs...)
		*pErr = zErrs.process(opt)
	}
}

func Appendf(err0 error, err error, msg string, args ...any) error {
	if err == nil {
		return err0
	}
	if msg != "" {
		err = Wrapf(err, msg, args...)
	}
	appendErrs(Option{}, &err0, err)
	return err0
}

func AppendTof(pErr *error, err error, msg string, args ...any) {
	if err == nil {
		return
	}
	if msg != "" {
		err = Wrapf(err, msg, args...)
	}
	appendErrs(Option{}, pErr, err)
}

func Validate(condition bool, msgArgs ...any) error {
	if condition {
		return nil
	}
	return New(formatValidate(msgArgs))
}

func Validatef(condition bool, msg string, args ...any) error {
	if condition {
		return nil
	}
	return Newf(msg, args...)
}

func ValidateX[T any](value T, condition bool, msgArgs ...any) (T, error) {
	if condition {
		return value, nil
	}
	return value, New(formatValidate(msgArgs))
}

func ValidateXf[T any](value T, condition bool, msgArgs ...any) (T, error) {
	if condition {
		return value, nil
	}
	return value, New(formatValidate(msgArgs))
}

func MustValidate(condition bool, msgArgs ...any) {
	if condition {
		return
	}
	panic(formatValidate(msgArgs))
}

func MustValidatef(condition bool, msg string, args ...any) {
	if condition {
		return
	}
	panic(sprintf(msg, args...))
}

func MustValidateX[T any](value T, condition bool, msgArgs ...any) T {
	if condition {
		return value
	}
	panic(formatValidate(msgArgs))
}

func MustValidateXf[T any](value T, condition bool, msg string, args ...any) T {
	if condition {
		return value
	}
	panic(sprintf(msg, args...))
}

func ValidateTo(pErr *error, condition bool, msgArgs ...any) {
	if !condition {
		msg := formatValidate(msgArgs)
		appendErrs(Option{}, pErr, NoStack().Error(msg))
	}
}

func ValidateTof(pErr *error, condition bool, msg string, args ...any) {
	if !condition {
		appendErrs(Option{}, pErr, NoStack().Errorf(msg, args...))
	}
}

func ValidateToX[T any](pErr *error, value T, condition bool, msgArgs ...any) (out T) {
	if condition {
		return value
	}
	msg := formatValidate(msgArgs)
	appendErrs(Option{}, pErr, NoStack().Error(msg))
	return out
}

func ValidateToXf[T any](pErr *error, value T, condition bool, msg string, args ...any) (out T) {
	if condition {
		return value
	}
	appendErrs(Option{}, pErr, NoStack().Errorf(msg, args...))
	return out
}

func GetErrors(err error) []error {
	switch err := err.(type) {
	case interface{ Errors() []error }:
		if err != nil {
			return err.Errors()
		}
	case interface{ WrappedErrors() []error }:
		if err != nil {
			return err.WrappedErrors()
		}
	}
	return nil
}

type zErrors struct {
	errors []error // len(errors) > 0
	stack  *stacktracez.Frames
}

func (es *zErrors) Error() string {
	switch {
	case es == nil:
		return "<nil>"
	case len(es.errors) == 0:
		return "<empty>"
	case len(es.errors) == 1:
		return sprintf("(1 error) %v", es.errors[0])
	default:
		return sprintf("%v", es)
	}
}

func (es *zErrors) Format(s0 fmt.State, verb rune) {
	s := fmtz.WrapState(s0)
	if es == nil {
		s.WriteStringZ("<nil>")
		return
	}
	isPlus := s.Flag('+') || s.Flag('#')
	switch verb {
	case 's', 'v':
		switch {
		case len(es.errors) == 0:
			s.WriteStringZ("<empty>")
			return
		case len(es.errors) == 1:
			if isPlus {
				s.Printf("1 error occurred:\n\t* %v\n", es.errors[0])
			} else {
				s.Printf("(1 error) %v", es.errors[0])
			}
			if isPlus && es.stack != nil {
				es.stack.Format(s, verb)
			}
		default:
			if isPlus {
				s.Printf("%d errors occurred:\n", len(es.errors))
			} else {
				s.Printf("(%d errors) ", len(es.errors))
			}
			for i, err := range es.errors {
				if isPlus {
					s.Printf("\t* %v\n", err)
					continue
				}
				if i > 0 {
					s.Printf(" ; ")
				}
				s.Printf("%v", err)
			}
			if isPlus && es.stack != nil {
				es.stack.Format(s, verb)
			}
		}
	case 'd':
		s.Printf("%d", len(es.errors))
	}
}

func (es *zErrors) Errors() []error {
	if es == nil {
		return nil
	}
	return es.errors[:len(es.errors)]
}

func (es *zErrors) Unwrap() []error {
	if es == nil {
		return nil
	}
	return es.Errors()
}

func (es *zErrors) StackTraceZ() *stacktracez.Frames {
	if es == nil {
		return nil
	}
	return es.stack
}

func (es *zErrors) Append(errs ...error) {
	for _, err := range errs {
		switch err := err.(type) {
		case nil:
			// continue
		case *zErrors:
			if err != nil {
				es.errors = append(es.errors, err.errors...)
			}
		case interface{ Unwrap() []error }:
			if err != nil {
				es.errors = append(es.errors, err.Unwrap()...)
			}
		case interface{ Errors() []error }:
			// uber-go/multierr, tailscale.com/util/multierr
			if err != nil {
				es.errors = append(es.errors, err.Errors()...)
			}
		case interface{ WrappedErrors() []error }:
			// hashicorp/go-multierror
			if err != nil {
				es.errors = append(es.errors, err.WrappedErrors()...)
			}
		default:
			if !reflect.ValueOf(err).IsNil() {
				es.errors = append(es.errors, err)
			}
		}
	}
}

func (es *zErrors) Appendf(err error, msgArgs ...any) {
	err = formatErrorMsg(err, msgArgs)
	if err != nil {
		es.errors = append(es.errors, err)
	}
}

func (es *zErrors) Validatef(condition bool, msg string, args ...any) {
	if !condition {
		err := NoStack().Errorf(msg, args...)
		es.errors = append(es.errors, err)
	}
}

func (es *zErrors) process(opt Option) error {
	if len(es.errors) == 0 {
		return nil
	}
	if es.stack != nil || opt.NoStack {
		return es
	}
	es.stack = stacktracez.StackTraceSkip(opt.CallersSkip + 3)
	return es
}

func sprintf(msg string, args ...any) string {
	if len(args) == 0 {
		return msg
	}
	return fmt.Sprintf(msg, args...)
}

func formatErrorMsg(err error, msgArgs []any) error {
	msg := fmtz.FormatMsgArgs(msgArgs)
	if msg == "" {
		return err
	}
	return NoStack().Wrap(err, msg)
}

func formatValidate(msgArgs []any) string {
	msg := fmtz.FormatMsgArgs(msgArgs)
	if msg == "" {
		return "validation failed"
	}
	return msg
}
