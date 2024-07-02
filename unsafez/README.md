<div align="center">

[![gopherz](https://ezpkg.io/_/gopherz.svg)](https://ezpkg.io)

</div>

# ezpkg.io/unsafez

[![PkgGoDev](https://pkg.go.dev/badge/ezpkg.io/unsafez)](https://pkg.go.dev/ezpkg.io/unsafez)
[![GitHub License](https://img.shields.io/github/license/ezpkg/unsafez)](https://github.com/ezpkg/unsafez/tree/main/LICENSE)
[![version](https://img.shields.io/github/v/tag/ezpkg/unsafez?label=version)](https://pkg.go.dev/ezpkg.io/unsafez?tab=versions)

Package [unsafez](https://pkg.go.dev/ezpkg.io/unsafez) provides unsafe functions for working with unsafe.Pointer.

## Installation

```sh
go get -u ezpkg.io/unsafez@v0.0.9
```

## Examples

```go
data := errorz.Must(os.ReadFile()))
str := unsafez.BytesToString(data)
fmt.Println(str)
```

## About ezpkg.io

As I work on various Go projects, I often find myself creating utility functions, extending existing packages, or developing packages to solve specific problems. Moving from one project to another, I usually have to copy or rewrite these solutions. So I created this repository to have all these utilities and packages in one place. Hopefully, you'll find them useful as well.

For more information, see the [main repository](https://github.com/ezpkg/ezpkg).

## Author

[![Oliver Nguyen](https://olivernguyen.io/_/badge.svg)](https://olivernguyen.io)&nbsp;&nbsp;[![github](https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white)](https://github.com/iOliverNguyen)
