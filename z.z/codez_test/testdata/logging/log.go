package logging

import (
	"fmt"

	"golang.org/x/net/context"
)

type Logger struct {
	ctx context.Context
}

func NewLogger(ctx context.Context) *Logger {
	return &Logger{ctx: ctx}
}

func (l *Logger) Log(msg string, args ...any) {
	fmt.Printf(msg, args...)
	fmt.Println()
}

func (l *Logger) LogX(ctx context.Context, msg string, args ...any) {
	fmt.Printf(msg, args...)
	fmt.Println()
}
