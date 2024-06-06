package errorz

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ezpkg/stacktracez"
)

var _ Errors = (*zErrors)(nil)
var _ stacktracez.StackTracerZ = (*zErrors)(nil)

type Errors interface {
	error
	Errors() []error
}

func Append(pErr *error, errs ...error) {
	switch err0 := (*pErr).(type) {
	case nil:
		var zErrs zErrors
		zErrs.Append(errs...)
		*pErr = zErrs.process()

	case *zErrors:
		err0.Append(errs...)
		*pErr = err0.process()

	default:
		var zErrs zErrors
		zErrs.errors = make([]error, len(errs))
		zErrs.Append(errs...)
		*pErr = zErrs.process()
	}
}

func Appendf(pErr *error, err error, msg string, args ...any) {
	if err == nil {
		return
	}
	if msg != "" {
		err = Wrapf(err, msg, args...)
	}
	Append(pErr, err)
}

func Validatef(pErr *error, condition bool, msg string, args ...any) {
	if !condition {
		Append(pErr, NoStack().Errorf(msg, args...))
	}
}

func ValidateX[T any](pErr *error, value T, condition bool, msg string, args ...any) (out T) {
	if condition {
		return value
	} else {
		Append(pErr, NoStack().Errorf(msg, args...))
		return out
	}
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
		return fmt.Sprintf("(1 error) %v", es.errors[0])
	default:
		return fmt.Sprintf("%v", es)
	}
}

func (es *zErrors) Format(s fmt.State, verb rune) {
	if es == nil {
		writeString(s, "<nil>")
		return
	}
	isPlus := s.Flag('+') || s.Flag('#')
	switch verb {
	case 's', 'v':
		switch {
		case len(es.errors) == 0:
			writeString(s, "<empty>")
			return
		case len(es.errors) == 1:
			if isPlus {
				fprintf(s, "1 error occurred:\n\t* %v\n", es.errors[0])
			} else {
				fprintf(s, "(1 error) %v", es.errors[0])
			}
			if isPlus && es.stack != nil {
				es.stack.Format(s, verb)
			}
		default:
			if isPlus {
				fprintf(s, "%d errors occurred:\n", len(es.errors))
			} else {
				fprintf(s, "(%d errors) ", len(es.errors))
			}
			for i, err := range es.errors {
				if isPlus {
					fprintf(s, "\t* %v\n", err)
					continue
				}
				if i > 0 {
					fprintf(s, " ; ")
				}
				fprintf(s, "%v", err)
			}
			if isPlus && es.stack != nil {
				es.stack.Format(s, verb)
			}
		}
	case 'd':
		fprintf(s, "%d", len(es.errors))
	}
}

func (es *zErrors) Errors() []error {
	if es == nil {
		return nil
	}
	return es.errors[:len(es.errors)]
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
	if err == nil {
		return
	}
	msg := formatMsg(msgArgs)
	if msg != "" {
		err = Wrap(err, msg)
	}
	es.errors = append(es.errors, err)
}

func (es *zErrors) Validatef(condition bool, msg string, args ...any) {
	if !condition {
		err := NoStack().Errorf(msg, args...)
		es.errors = append(es.errors, err)
	}
}

func (es *zErrors) process() error {
	if len(es.errors) == 0 {
		return nil
	}
	if es.stack == nil {
		es.stack = stacktracez.StackTraceSkip(2)
	}
	return es
}

func formatMsg(msgArgs []any) string {
	if len(msgArgs) == 0 {
		return ""
	}
	if format, ok := msgArgs[0].(string); ok {
		return fmt.Sprintf(format, msgArgs[1:]...)
	}
	return strings.TrimSpace(fmt.Sprintln(msgArgs...))
}
