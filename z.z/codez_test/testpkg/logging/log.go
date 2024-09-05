package logging

import (
	"context"
	"errors"
	"fmt"
)

type Logger struct {
	ctx context.Context
}

func NewLogger(ctx context.Context) (*Logger, error) {
	return &Logger{ctx: ctx}, nil
}

func (l *Logger) Log(msg string, args ...any) {
	fmt.Printf(msg, args...)
	fmt.Println()
}
func (l *Logger) LogX(ctx context.Context, msg string, args ...any) {
	fmt.Printf(msg, args...)
	fmt.Println()
}
func (l *Logger) GetError() error {
	return nil
}
func (l *Logger) AllErrors() []error {
	return nil
}
func (l *Logger) TryOne() (string, error) {
	return "", nil
}
func (l *Logger) TryTwo() (a, b string, err error) {
	return "", "", nil
}

func ErrorOne() error {
	return errors.New("one")
}
