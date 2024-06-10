<p align="center">
<a href="https://ezpkg.io">
<img alt="gopherz" src="https://ezpkg.io/_/gopherz.png" style="width:420px">
</a>
</p>

# ezpkg.io/stacktracez

[![PkgGoDev](https://pkg.go.dev/badge/github.com/ezpkg/stacktracez)](https://pkg.go.dev/ezpkg.io/stacktracez)
[![GitHub License](https://img.shields.io/github/license/ezpkg/stacktracez)](https://github.com/ezpkg/stacktracez/tree/main/LICENSE)
[![version](https://img.shields.io/github/v/tag/ezpkg/stacktracez?label=version)](https://github.com/ezpkg/stacktracez/tags)

Package stacktracez provides functions for getting stack trace for using in errors and logs.

## Installation

```sh
go get -u ezpkg.io/stacktracez@v0.0.3
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

## About ezpkg.io

As I work on various Go projects, I often find myself creating utility functions, extending existing packages, or developing packages to solve specific problems. Moving from one project to another, I usually have to copy or rewrite these solutions. So I created this repository to have all these utilities and packages in one place. Hopefully, you'll find them useful as well.

For more information, see the [main repository](https://github.com/ezpkg/ezpkg).

## Author

<a href="https://olivernguyen.io"><img alt="olivernguyen.io" src="https://olivernguyen.io/_/badge.png" height="28px"></a>&nbsp;&nbsp;[![github](https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white)](https://github.com/iOliverNguyen)
