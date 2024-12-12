<div align="center">

[![gopherz](https://ezpkg.io/_/gopherz.svg)](https://ezpkg.io)

</div>

# ezpkg.io/diffz

[![PkgGoDev](https://pkg.go.dev/badge/ezpkg.io/diffz)](https://pkg.go.dev/ezpkg.io/diffz)
[![GitHub License](https://img.shields.io/github/license/ezpkg/diffz)](https://github.com/ezpkg/diffz/tree/main/LICENSE)
[![version](https://img.shields.io/github/v/tag/ezpkg/diffz?label=version)](https://pkg.go.dev/ezpkg.io/diffz?tab=versions)

Package [diffz](https://pkg.go.dev/ezpkg.io/diffz) provides functions for comparing and displaying differences between two strings. It's based on [kylelemons/godebug](https://github.com/kylelemons/godebug) and [sergi/go-diff](https://github.com/sergi/go-diff). It provides additional features of ignoring spaces and supporting placeholders.

## Installation

```sh
go get -u ezpkg.io/diffz@v0.2.1
```

## Examples

```go
// diff by char
left, right := "onetwo thr33four five", "onetwothree fourfive"
diffs := diffz.IgnoreSpace().DiffByChar(left, right)
fmt.Println(diffz.Format(diffs))

// diff by line
left, right := "one\ntwo\nthree\nfour", "one\ntwo\nfour"
diffs := diffz.IgnoreSpace().DiffByLine(left, right)
fmt.Println(diffz.Format(diffs))

// placeholder is useful for comparing tests with uuid or random values
diffs := diffz.Placeholder().AndIgnoreSpace().DiffByLine(left, right)
left := "id: ████\ncode: AF███\nname: Alice\n"
right := "id: 1234\ncode: AF123\nname: Baby\n"
fmt.Println(diffz.Format(diffs))
```

## Similar Packages

This package is based on these packages:

- [github.com/kylelemons/godebug](https://github.com/kylelemons/godebug)
- [github.com/sergi/go-diff](https://github.com/sergi/go-diff)

## About ezpkg.io

As I work on various Go projects, I often find myself creating utility functions, extending existing packages, or developing packages to solve specific problems. Moving from one project to another, I usually have to copy or rewrite these solutions. So I created this repository to have all these utilities and packages in one place. Hopefully, you'll find them useful as well.

For more information, see the [main repository](https://github.com/ezpkg/ezpkg).

## Author

[![Oliver Nguyen](https://olivernguyen.io/_/badge.svg)](https://olivernguyen.io)&nbsp;&nbsp;[![github](https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white)](https://github.com/iOliverNguyen)
