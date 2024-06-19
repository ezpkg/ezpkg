<div align="center">

[![gopherz](https://ezpkg.io/_/gopherz.svg)](https://ezpkg.io)

</div>

# ezpkg.io/stacktracez

[![PkgGoDev](https://pkg.go.dev/badge/ezpkg.io/stacktracez)](https://pkg.go.dev/ezpkg.io/stacktracez)
[![GitHub License](https://img.shields.io/github/license/ezpkg/stacktracez)](https://github.com/ezpkg/stacktracez/tree/main/LICENSE)
[![version](https://img.shields.io/github/v/tag/ezpkg/stacktracez?label=version)](https://github.com/ezpkg/stacktracez/tags)

Package [stacktracez](https://pkg.go.dev/ezpkg.io/stacktracez) provides functions for getting stack trace for using in errors and logs.

## Installation

```sh
go get -u ezpkg.io/stacktracez@v0.0.7
```

## Examples

```go
stack := stacktracez.StackTrace()

fmt.Printf("%+v", stack)
// ezpkg.io/stacktracez.TestStackTrace.func1
// /Users/i/ws/ezpkg/ezpkg/stacktracez/stacktracez_test.go:12
// ezpkg.io/stacktracez.TestStackTrace
// /Users/i/ws/ezpkg/ezpkg/stacktracez/stacktracez_test.go:15
// testing.tRunner
// /usr/local/go/src/testing/testing.go:1689

fmt.Printf("%v", stack)
// ezpkg.io/stacktracez/stacktracez_test.go:12 · TestStackTrace.func1
// ezpkg.io/stacktracez/stacktracez_test.go:15 · TestStackTrace
// testing/testing.go:1689 · tRunner
```

## Similar Packages

- [github.com/pkg/errors](https://github.com/pkg/errors)
- [go.elastic.co/apm/v2/stacktrace](https://pkg.go.dev/go.elastic.co/apm/v2/stacktrace)
- [github.com/palantir/stacktrace](github.com/palantir/stacktrace)

## About ezpkg.io

As I work on various Go projects, I often find myself creating utility functions, extending existing packages, or developing packages to solve specific problems. Moving from one project to another, I usually have to copy or rewrite these solutions. So I created this repository to have all these utilities and packages in one place. Hopefully, you'll find them useful as well.

For more information, see the [main repository](https://github.com/ezpkg/ezpkg).

## Author

[![Oliver Nguyen](https://olivernguyen.io/_/badge.svg)](https://olivernguyen.io)&nbsp;&nbsp;[![github](https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white)](https://github.com/iOliverNguyen)
