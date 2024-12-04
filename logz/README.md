<div align="center">

[![gopherz](https://ezpkg.io/_/gopherz.svg)](https://ezpkg.io)

</div>

# ezpkg.io/logz

[![PkgGoDev](https://pkg.go.dev/badge/ezpkg.io/logz)](https://pkg.go.dev/ezpkg.io/logz)
[![GitHub License](https://img.shields.io/github/license/ezpkg/logz)](https://github.com/ezpkg/logz/tree/main/LICENSE)
[![version](https://img.shields.io/github/v/tag/ezpkg/logz?label=version)](https://pkg.go.dev/ezpkg.io/logz?tab=versions)

Package [logz](https://pkg.go.dev/ezpkg.io/logz) provides common interfaces and utilities for working with other log packages, including [log/slog](https://pkg.go.dev/log/slog) and [zap](https://pkg.go.dev/go.uber.org/zap). It's useful for package library to provide a common interface for logging, so that the user can choose the logging package they want to use.

## Installation

```sh
go get -u ezpkg.io/logz@v0.2.0
```

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

## About ezpkg.io

As I work on various Go projects, I often find myself creating utility functions, extending existing packages, or developing packages to solve specific problems. Moving from one project to another, I usually have to copy or rewrite these solutions. So I created this repository to have all these utilities and packages in one place. Hopefully, you'll find them useful as well.

For more information, see the [main repository](https://github.com/ezpkg/ezpkg).

## Author

[![Oliver Nguyen](https://olivernguyen.io/_/badge.svg)](https://olivernguyen.io)&nbsp;&nbsp;[![github](https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white)](https://github.com/iOliverNguyen)
