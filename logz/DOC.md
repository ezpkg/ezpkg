# ezpkg.io/logz

Package [logz](https://pkg.go.dev/ezpkg.io/logz) provides common interfaces and utilities for working with other log packages, including [log/slog](https://pkg.go.dev/log/slog) and [zap](https://pkg.go.dev/go.uber.org/zap). It's useful for package library to provide a common interface for logging, so that the user can choose the logging package they want to use.

## Examples

```go
package mypackage

import "log/slog"
import "os"
import "testing"
import "go.uber.org/zap"

import "ezpkg.io/errorz"
import "ezpkg.io/logz"

type SuperPower struct {
	logger logz.Logger
}

func NewSuperPower(logger logz.Logger) (*SuperPower, error) {
	sp := &SuperPower{logger: logger}
	sp.logger.Debugw("A new SuperPower is created! Better be prepared! 🔥")
	return sp, nil
}

func TestSuperPowerUsingSlog(t *testing.T) {
	opt := &slog.HandlerOptions{Level: slog.LevelDebug}
	handler := slog.NewTextHandler(os.Stderr, opt)
	stdLogger := slog.New(handler)

	sp, err := NewSuperPower(logz.FromLoggerI(stdLogger))
	assertSuperPower(sp, err)
}
func TestSuperPowerUsingZap(t *testing.T) {
	zapLogger := errorz.Must(zap.NewDevelopment()).Sugar()

	sp, err := NewSuperPower(logz.FromLoggerx(zapLogger))
	assertSuperPower(sp, err)
}
```
