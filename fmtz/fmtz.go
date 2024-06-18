package fmtz

import (
	"fmt"
	"io"
	"strings"
)

func FormatMsgArgs(msgAndArgs []any) string {
	if len(msgAndArgs) == 0 {
		return ""
	}
	if msg, ok := msgAndArgs[0].(string); ok {
		if len(msgAndArgs) == 1 {
			return msg
		}
		return fmt.Sprintf(msg, msgAndArgs[1:]...)
	}
	return strings.TrimSpace(fmt.Sprintln(msgAndArgs...))
}

func FormatMsgArgsX(msgAndArgs []any) fmt.Formatter {
	if len(msgAndArgs) == 0 {
		return emptyFormatMsgArgsX
	}
	return implFormatMsgArgsX{msgAndArgs: msgAndArgs}
}

var emptyFormatMsgArgsX = implFormatMsgArgsX{}

type implFormatMsgArgsX struct {
	msgAndArgs []any
}

func (x implFormatMsgArgsX) Format(s fmt.State, verb rune) {
	if len(x.msgAndArgs) == 0 {
		return
	}
	if msg, ok := x.msgAndArgs[0].(string); ok {
		if len(x.msgAndArgs) == 1 {
			_, _ = io.WriteString(s, msg)
		}
		_, _ = fmt.Fprintf(s, msg, x.msgAndArgs[1:]...)
	}
	for i, v := range x.msgAndArgs {
		if i > 0 {
			_, _ = io.WriteString(s, " ")
		}
		_, _ = fmt.Fprint(s, v)
	}
}
