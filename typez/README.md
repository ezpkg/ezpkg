<div align="center">

[![gopherz](https://ezpkg.io/_/gopherz.svg)](https://ezpkg.io)

</div>

# ezpkg.io/typez

[![PkgGoDev](https://pkg.go.dev/badge/ezpkg.io/typez)](https://pkg.go.dev/ezpkg.io/typez)
[![GitHub License](https://img.shields.io/github/license/ezpkg/typez)](https://github.com/ezpkg/typez/tree/main/LICENSE)
[![version](https://img.shields.io/github/v/tag/ezpkg/typez?label=version)](https://pkg.go.dev/ezpkg.io/typez?tab=versions)

Package [typez](https://pkg.go.dev/ezpkg.io/typez) provides generic functions for working with types.

## Installation

```sh
go get -u ezpkg.io/typez@v0.0.9
```

## Examples

```go
typez.In(1, 1, 2, 3)    // true
typez.In("A", "B", "C") // false

type A struct{X int}
typez.Coalesce(0, 1, 2, 3) // 1
typez.Coalesce(nil, &A{10}, &A{20}) // &A{10}
```

## About ezpkg.io

As I work on various Go projects, I often find myself creating utility functions, extending existing packages, or developing packages to solve specific problems. Moving from one project to another, I usually have to copy or rewrite these solutions. So I created this repository to have all these utilities and packages in one place. Hopefully, you'll find them useful as well.

For more information, see the [main repository](https://github.com/ezpkg/ezpkg).

## Author

[![Oliver Nguyen](https://olivernguyen.io/_/badge.svg)](https://olivernguyen.io)&nbsp;&nbsp;[![github](https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white)](https://github.com/iOliverNguyen)
